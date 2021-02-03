local go_runtime(version, arch) = {
  type: 'pod',
  arch: arch,
  containers: [
    { image: 'golang:' + version },
  ],
};

local task_build_go() = {
  name: 'build go',
  runtime: go_runtime('1.15', 'amd64'),
  steps: [
    { type: 'clone' },
    { type: 'restore_cache', keys: ['cache-sum-{{ md5sum "go.sum" }}', 'cache-date-'], dest_dir: '/go/pkg/mod/cache' },
    {
      type: 'run',
      name: 'build',
      command: |||
        if [ ${AGOLA_GIT_TAG} ];
          then export SERVER_VERSION=${AGOLA_GIT_TAG} ;
        else
          export SERVER_VERSION=${AGOLA_GIT_COMMITSHA} ; fi

        echo ${SERVER_VERSION}

        go build -ldflags "-X github.com/ercole-io/ercole/v2/cmd.serverVersion=${SERVER_VERSION}"
      |||,
    },
    { type: 'save_to_workspace', contents: [{ source_dir: '.', dest_dir: '.', paths: ['ercole', 'package/**', 'resources/**', 'distributed_files/**'] }] },
  ],
  depends: ['test'],
};

local task_pkg_build(setup) = {
  name: 'pkg build ' + setup.dist,
  runtime: {
    type: 'pod',
    arch: 'amd64',
    containers: [
      { image: setup.pkg_build_image },
    ],
  },
  environment: {
    WORKSPACE: '/root/project',
    DIST: setup.dist,
  },
  steps: [
    { type: 'restore_workspace', dest_dir: '.' },
    { type: 'run', command: 'mkdir -p ${WORKSPACE}/dist' },
    { type: 'run', command: 'mkdir -p ~/rpmbuild/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}' },
    {
      type: 'run',
      name: 'ln & cd',
      command: |||
        if [ ${AGOLA_GIT_TAG} ];
          then export VERSION=${AGOLA_GIT_TAG} ;
        else
          export VERSION=latest ; fi

        echo ${VERSION}

        ln -s ${WORKSPACE} ~/rpmbuild/SOURCES/ercole-${VERSION}
        cd ${WORKSPACE} && rpmbuild --define "_version ${VERSION}" -bb package/${DIST}/ercole.spec
      |||,
    },
    { type: 'run', command: 'ls ~/rpmbuild/RPMS/x86_64/ercole-*.rpm' },
    { type: 'run', command: 'file ~/rpmbuild/RPMS/x86_64/ercole-*.rpm' },
    { type: 'run', command: 'cp ~/rpmbuild/RPMS/x86_64/ercole-*.rpm ${WORKSPACE}/dist' },
    { type: 'save_to_workspace', contents: [{ source_dir: './dist/', dest_dir: '/dist/', paths: ['**'] }] },
  ],
  depends: ['build go'],
};

local task_deploy_repository(dist) = {
  name: 'deploy repository.ercole.io ' + dist,
  approval: true,
  runtime: {
    type: 'pod',
    arch: 'amd64',
    containers: [
      { image: 'curlimages/curl' },
    ],
  },
  environment: {
    REPO_USER: { from_variable: 'repo-user' },
    REPO_TOKEN: { from_variable: 'repo-token' },
    REPO_UPLOAD_URL: { from_variable: 'repo-upload-url' },
    REPO_INSTALL_URL: { from_variable: 'repo-install-url' },
  },
  steps: [
    { type: 'restore_workspace', dest_dir: '.' },
    {
      type: 'run',
      name: 'curl',
      command: |||
        cd dist
        for f in *; do
        	URL=$(curl --user "${REPO_USER}" \
            --upload-file $f ${REPO_UPLOAD_URL} --insecure)
        	echo $URL
        	md5sum $f
        	curl -H "X-API-Token: ${REPO_TOKEN}" \
          -H "Content-Type: application/json" --request POST --data "{ \"filename\": \"$f\", \"url\": \"$URL\" }" \
          ${REPO_INSTALL_URL} --insecure
        done
      |||,
    },
  ],
  depends: ['pkg build ' + dist],
  when: {
    branch: 'master',
    ref: {
      exclude: ['#/refs/pull/.*#'],
    },
  },
};

local task_build_push_image(push) =
  /*
   * Currently, kaniko, has some issues with multi stage builds where it removes
   * all the files in the container after every stage (excluding /kaniko) causing
   * file not found errors when doing COPY commands.
   * Workaround this buy putting all files inside /kaniko
   */
  local options = if !push then '--no-push' else '--destination sorintlab/ercole-services:$AGOLA_GIT_TAG';
  {
    name: 'build image' + if push then ' and push' else '',
    runtime: {
      arch: 'amd64',
      containers: [
        {
          image: 'gcr.io/kaniko-project/executor:debug-v0.11.0',
        },
      ],
    },
    environment: {
      DOCKERAUTH: { from_variable: 'dockerauth' },
    },
    shell: '/busybox/sh',
    working_dir: '/kaniko',
    steps: [
      { type: 'restore_workspace', dest_dir: '/kaniko/ercole' },
    ] + std.prune([
      if push then {
        type: 'run',
        name: 'generate docker auth',
        command: |||
          cat << EOF > /kaniko/.docker/config.json
          {
            "auths": {
              "https://index.docker.io/v1/": { "auth" : "$DOCKERAUTH" }
            }
          }
          EOF
        |||,
      },
    ]) + [
      { type: 'run', command: '/kaniko/executor --context=dir:///kaniko/ercole --dockerfile Dockerfile %s' % [options] },
    ],
    depends: ['checkout code'],
  };

{
  runs: [
    {
      name: 'ercole',
      tasks: [
        {
          name: 'test',
          runtime: {
            type: 'pod',
            arch: 'amd64',
            containers: [
              { image: 'golang:1.15' },
              { image: 'mongo:4' },
            ],
          },
          steps: [
            { type: 'clone' },
            { type: 'restore_cache', keys: ['cache-sum-{{ md5sum "go.sum" }}', 'cache-date-'], dest_dir: '/go/pkg/mod/cache' },

            { type: 'run', name: '', command: 'go get github.com/golang/mock/mockgen@v1.4.4' },
            { type: 'run', name: '', command: 'go generate -v ./...' },
            { type: 'run', name: '', command: 'go test -race -coverprofile=coverage.txt -covermode=atomic ./...' },

            { type: 'save_cache', key: 'cache-sum-{{ md5sum "go.sum" }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
            { type: 'save_cache', key: 'cache-date-{{ year }}-{{ month }}-{{ day }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
          ],
        },
      ] + [
        task_build_go(),
      ] + [
        {
          name: 'version',
          runtime: {
            type: 'pod',
            arch: 'amd64',
            containers: [
              { image: 'debian:buster' },
            ],
          },
          steps: [
            { type: 'restore_workspace', dest_dir: '.' },
            { type: 'run', command: './ercole version' },
          ],
          depends: ['build go'],
        },
      ] + [
        task_pkg_build(setup)
        for setup in [
          { pkg_build_image: 'amreo/rpmbuild-centos7', dist: 'rhel7' },
          { pkg_build_image: 'amreo/rpmbuild-centos8', dist: 'rhel8' },
        ]
        //TODO Publish assets to GitHub
      ] + [
        task_deploy_repository(dist)
        for dist in ['rhel7', 'rhel8']
      ] + [
        {
          name: 'checkout code',
          runtime: {
            type: 'pod',
            arch: 'amd64',
            containers: [
              { image: 'alpine/git' },
            ],
          },
          steps: [
            { type: 'clone' },
            { type: 'save_to_workspace', contents: [{ source_dir: '.', dest_dir: '.', paths: ['**'] }] },
          ],
          depends: ['test'],
        },
      ] + [
        task_build_push_image(false) + {
          when: {
            ref: '#refs/pull/\\d+/head#',
          },
        },
        task_build_push_image(true) + {
          when: {
            branch: 'master',
            tag: '#v.*#',
          },
        },
      ] + [
        {
          name: 'redeploy dev.ercole.io',
          runtime: {
            type: 'pod',
            arch: 'amd64',
            containers: [
              { image: 'curlimages/curl' },
            ],
          },
          environment: {
            REDEPLOY_URL: { from_variable: 'redeploy-url' },
          },
          steps: [
            {
              type: 'run',
              name: 'curl request',
              command: |||
                curl --location --request POST ${REDEPLOY_URL} \
                  --header 'Content-Type: application/json' \
                  --data-raw '{ "namespace": "ercole", "podname" : "ercole-services" }' \
                  --insecure
              |||,
            },
          ],
          depends: ['build image and push'],
        },
      ],
    },
  ],
}

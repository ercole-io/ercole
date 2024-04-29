local go_runtime(version, arch) = {
  type: 'pod',
  arch: arch,
  containers: [
    { image: 'golang:' + version },
  ],
};

local task_build_go_rhel7() = {
  name: 'build go rhel7',
  runtime: {
    type: 'pod',
    containers: [
      { image: 'fra.ocir.io/fremyxlx6yog/oraclelinux7:latest' },
    ],
  },
  steps: [
    { type: 'run', command: 'git clone https://github.com/ercole-io/ercole.git .' },
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

local task_build_go_rhel8() = {
  name: 'build go rhel8',
  runtime: {
    type: 'pod',
    containers: [
      { image: 'fra.ocir.io/fremyxlx6yog/oraclelinux8:latest' },
    ],
  },
  steps: [
    { type: 'run', command: 'git clone https://github.com/ercole-io/ercole.git .' },
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

local task_build_go_rhel9() = {
  name: 'build go rhel9',
  runtime: {
    type: 'pod',
    containers: [
      { image: 'fra.ocir.io/fremyxlx6yog/oraclelinux9:latest' },
    ],
  },
  steps: [
    { type: 'run', command: 'git clone https://github.com/ercole-io/ercole.git .' },
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

local version_rhel8() = {
  name: 'version rhel8',
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
  depends: ['build go rhel8'],
};

local version_rhel7() = {
  name: 'version rhel7',
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
  depends: ['build go rhel7'],
};

local version_rhel9() = {
  name: 'version rhel9',
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
  depends: ['build go rhel9'],
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

        SUPPORTED_VERSION=$(echo $VERSION | sed 's/-/_/g')

        echo ${SUPPORTED_VERSION}

        ln -s ${WORKSPACE} ~/rpmbuild/SOURCES/ercole-${SUPPORTED_VERSION}
        cd ${WORKSPACE} && rpmbuild --define "_version ${SUPPORTED_VERSION}" -bb package/${DIST}/ercole.spec
      |||,
    },
    { type: 'run', command: 'ls ~/rpmbuild/RPMS/x86_64/ercole-*.rpm' },
    { type: 'run', command: 'file ~/rpmbuild/RPMS/x86_64/ercole-*.rpm' },
    { type: 'run', command: 'cp ~/rpmbuild/RPMS/x86_64/ercole-*.rpm ${WORKSPACE}/dist' },
    { type: 'save_to_workspace', contents: [{ source_dir: './dist/', dest_dir: '/dist/', paths: ['**'] }] },
  ],
  depends: ['build go '+ setup.dist],
};

local task_deploy_repository(dist) = {
  name: 'upload to repository.ercole.io ' + dist,
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
      name: 'upload to repository.ercole.io',
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
    tag: '#.*#',
    branch: 'master',
  },
};

local task_upload_asset(dist) = {
 name: 'upload to github.com ' + dist,
  runtime: {
    type: 'pod',
    arch: 'amd64',
    containers: [
      { image: 'curlimages/curl' },
    ],
  },
 environment: {
    GITHUB_USER: { from_variable: 'github-user' },
    GITHUB_TOKEN: { from_variable: 'github-token' },
  },
steps: [
    { type: 'restore_workspace', dest_dir: '.' },
    {
      type: 'run',
      name: 'upload to github',
      command: |||
          cd dist
          GH_REPO="https://api.github.com/repos/${GITHUB_USER}/ercole/releases"
          if [ ${AGOLA_GIT_TAG} ];
            then
              GH_TAGS="$GH_REPO/tags/$AGOLA_GIT_TAG" ;
              response=$(curl -sH "Authorization: token ${GITHUB_TOKEN}" $GH_TAGS) ;
              eval $(echo "$response" | grep -m 1 "id.:" | grep -w id | tr : = | tr -cd '[[:alnum:]]=') ; 
              for filename in *; do
                REPO_ASSET="https://uploads.github.com/repos/${GITHUB_USER}/ercole/releases/$id/assets?name=$(basename $filename)" ;
                curl -H POST -H "Authorization: token ${GITHUB_TOKEN}" -H "Content-Type: application/octet-stream" --data-binary @"$filename" $REPO_ASSET ;
                echo $REPO_ASSET ;
              done
          fi
      |||,
    },
  ],
  depends: ['pkg build ' + dist],
  when: {
    tag: '#.*#',
    branch: 'master',
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
              { image: 'golang:1.20' },
              { image: 'mongo:4' },
            ],
          },
          environment: {
            GITLEAKS_CONF: { from_variable: 'gitleaks-config' },
          },
          steps: [
            { type: 'clone' },
            { type: 'restore_cache', keys: ['cache-sum-1-{{ md5sum "go.sum" }}', 'cache-date-1-'], dest_dir: '/go/pkg/mod/cache' },

            { type: 'run', name: 'clone gitleaks', command: 'git clone https://github.com/gitleaks/gitleaks.git ../gitleaks' },
            { type: 'run', name: 'build gitleaks', command: 'cd ../gitleaks; make build; echo ${GITLEAKS_CONF} > gitleaks.toml' },
            { type: 'run', name: 'detect security leaks', command: '../gitleaks/gitleaks detect -v -c ../gitleaks/gitleaks.toml' },
  
            { type: 'run', name: 'install golangci-lint', command: 'curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2' },
            { type: 'run', name: 'run golangci-lint', command: 'golangci-lint run' },

            { type: 'run', name: '', command: 'go install go.uber.org/mock/mockgen@v0.3.0' },
            { type: 'run', name: '', command: 'go generate -v ./...' },
            { type: 'run', name: '', command: 'go test -race -coverprofile=coverage.txt -covermode=atomic ./...' },

            { type: 'save_cache', key: 'cache-sum-1-{{ md5sum "go.sum" }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
            { type: 'save_cache', key: 'cache-date-1-{{ year }}-{{ month }}-{{ day }}', contents: [{ source_dir: '/go/pkg/mod/cache' }] },
          ],
        },
      ] + [
        task_build_go_rhel7(),
      ] + [
        task_build_go_rhel8(),
      ] + [
        task_build_go_rhel9(),
      ] + [
        version_rhel7()
      ] + [
        version_rhel8()
      ] + [
        version_rhel9()
      ] + [
        task_pkg_build(setup)
        for setup in [
          { pkg_build_image: 'amreo/rpmbuild-centos7', dist: 'rhel7' },
          { pkg_build_image: 'amreo/rpmbuild-centos8', dist: 'rhel8' },
          { pkg_build_image: 'fra.ocir.io/fremyxlx6yog/rpmbuildrhel9', dist: 'rhel9' },
        ]       
      ] + [  
        task_deploy_repository(dist)
        for dist in ['rhel7', 'rhel8', 'rhel9']
      ] + [
        task_upload_asset(dist)
        for dist in ['rhel7', 'rhel8', 'rhel9']
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
            branch: '^(?!master$).*$',
            ref: '#refs/pull/\\d+/head#',
          },
        },
        task_build_push_image(true) + {
          when: {
            branch: 'master',
            tag: '#.*#',
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
                  --data-raw '{ "namespace": "default", "podname" : "ercole-services" }' \
                  --insecure
              |||,
            },
          ],
          depends: ['build image and push'],
        },
      ],
      docker_registries_auth: {
        'index.docker.io': {
	  type: 'basic',
	  username: { from_variable: 'docker-username' },
	  password: { from_variable: 'docker-password' },
	},
      },
    },
  ],
}
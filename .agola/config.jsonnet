local task_build_push_image() =
  /*
   * Currently, kaniko, has some issues with multi stage builds where it removes
   * all the files in the container after every stage (excluding /kaniko) causing
   * file not found errors when doing COPY commands.
   * Workaround this buy putting all files inside /kaniko
   */
  {
    name: 'build & push image',
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
      {
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
      { type: 'run', command: '/kaniko/executor --context=dir:///kaniko/ercole --dockerfile Dockerfile --destination sorintlab/ercole-services:$AGOLA_GIT_TAG' },
    ],
    depends: ['checkout code'],
  };

local task_build_go_rhel(rhel_version) = {
  name: 'build go rhel' + rhel_version,
  runtime: {
    type: 'pod',
    containers: [
      { image: 'fra.ocir.io/fremyxlx6yog/oraclelinux' + rhel_version + ':latest' },
    ],
  },
  steps: [
    { type: 'clone' },
    {
      type: 'run',
      name: 'build',
      command: |||
        if [ ${AGOLA_GIT_TAG} ];
          then export SERVER_VERSION=${AGOLA_GIT_TAG} ;
        else
          export SERVER_VERSION=${AGOLA_GIT_COMMITSHA} ;
        fi
        go build -ldflags "-X github.com/ercole-io/ercole/v2/cmd.serverVersion=${SERVER_VERSION}"
      |||,
    },
    { type: 'save_to_workspace', contents: [{ source_dir: '.', dest_dir: '.', paths: ['ercole', 'package/**', 'resources/**', 'distributed_files/**'] }] },
  ],
  depends: ['check secrets'],
};

local version_rhel(rhel_version) = {
  name: 'version rhel' + rhel_version,
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
  depends: ['build go rhel' + rhel_version],
};

local task_pkg_build(setup) = {
  name: 'pkg build rhel' + setup.dist,
  runtime: {
    type: 'pod',
    arch: 'amd64',
    containers: [
      { image: setup.pkg_build_image },
    ],
  },
  environment: {
    WORKSPACE: '/root/project',
    DIST: 'rhel' + setup.dist,
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

        SUPPORTED_VERSION=$(echo $VERSION | sed 's/-/_/g')

        ln -s ${WORKSPACE} ~/rpmbuild/SOURCES/ercole-${SUPPORTED_VERSION}
        cd ${WORKSPACE} && rpmbuild --define "_version ${SUPPORTED_VERSION}" -bb package/${DIST}/ercole.spec
      |||,
    },
    { type: 'run', command: 'ls ~/rpmbuild/RPMS/x86_64/ercole-*.rpm' },
    { type: 'run', command: 'file ~/rpmbuild/RPMS/x86_64/ercole-*.rpm' },
    { type: 'run', command: 'cp ~/rpmbuild/RPMS/x86_64/ercole-*.rpm ${WORKSPACE}/dist' },
    { type: 'save_to_workspace', contents: [{ source_dir: './dist/', dest_dir: '/dist/', paths: ['**'] }] },
  ],
  depends: [ 'version rhel' + setup.dist ],
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

{
  docker_registries_auth: {
    'index.docker.io': {
      type: 'basic',
      username: { from_variable: 'docker-username' },
      password: { from_variable: 'docker-password' },
    },
  },
  runs: [ 
    {
      name: 'ercole',
      tasks: [
        {
          name: 'linters',
          runtime: {
            type: 'pod',
            arch: 'amd64',
            containers: [
              { image: 'golangci/golangci-lint:v1.52.2' },
            ],
          },
          steps: [
            { type: 'clone' },
            { type: 'run', name: 'run golangci-lint', command: 'golangci-lint run --timeout 10m' },
          ],
        },
      ] + [
        {
          name: 'staticcheck',
          runtime: {
            type: 'pod',
            arch: 'amd64',
            containers: [
              { image: 'golang:1.22' },
            ],
          },
          steps: [
            { type: 'clone' },
            { type: 'run', name: 'install staticcheck', command: 'go install honnef.co/go/tools/cmd/staticcheck@latest' },
            { type: 'run', name: 'run staticcheck', command: 'staticcheck -f=stylish -tests=false ./...' },
          ],
          depends: ['linters'],
        },
      ] + [
        {
          name: 'test',
          runtime: {
            type: 'pod',
            arch: 'amd64',
            containers: [
              { image: 'golang:1.21' },
              { image: 'mongo:6' },
            ],
          },
          steps: [
            { type: 'clone' },
            { type: 'run', name: '', command: 'go install go.uber.org/mock/mockgen@v0.3.0' },
            { type: 'run', name: '', command: 'go generate -v ./...' },
            { type: 'run', name: '', command: 'go test -race ./... -v' },
          ],
          depends: ['staticcheck'],
        }
      ] + [
        {
          name: 'check secrets',
          runtime: {
            type: 'pod',
            arch: 'amd64',
            containers: [
              { image: 'zricethezav/gitleaks' },
            ],
          },
          environment: {
            GITLEAKS_CONF: { from_variable: 'gitleaks-config' },
          },
          steps: [
            { type: 'clone' },
            { type: 'run', name: 'write gitleaks config', command: 'printf "%b" ${GITLEAKS_CONF} > gitleaks.toml' },
            { type: 'run', name: 'detect security leaks', command: 'gitleaks detect -v -c gitleaks.toml' },
          ],
          depends: [
            'test'
          ],
        },
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
          depends: [
            'check secrets'
          ],
          when: {
            branch: 'master',
            tag: '#.*#',
          },
        },
      ] + [
        task_build_push_image() + {
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
          depends: ['build & push image'],
          when: {
            tag: '#.*#',
            branch: 'master',
          },
        },
      ] + [
        task_build_go_rhel(rhel_version)
        for rhel_version in ['7', '8', '9']
      ] + [
        version_rhel(rhel_version)
        for rhel_version in ['7', '8', '9']
      ] + [
        task_pkg_build(setup)
        for setup in [
          { pkg_build_image: 'amreo/rpmbuild-centos7', dist: '7' },
          { pkg_build_image: 'amreo/rpmbuild-centos8', dist: '8' },
          { pkg_build_image: 'fra.ocir.io/fremyxlx6yog/rpmbuildrhel9', dist: '9' },
        ]
      ] + [  
        task_deploy_repository(dist)
        for dist in ['rhel7', 'rhel8', 'rhel9']
      ],
    },
  ],
}
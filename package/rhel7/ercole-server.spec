Name:           ercole-server
Version:        ERCOLE_VERSION
Release:        1%{?dist}
Summary:        Ercole server	

License:        ASL 2.0
URL:            https://ercole.io            
Source0:        https://github.com/ercole-io/%{name}/archive/%{version}.tar.gz
Group:          Tools
Requires:       java-11-openjdk systemd
BuildRequires: systemd

Buildroot:      /tmp/rpm-ercole-server
%global         debug_package %{nil}
%define         __jar_repack 0

%description
This is the server component for the Ercole project.

%global debug_package %{nil}

%pre
    getent passwd ercole >/dev/null || useradd -s /bin/bash -c "Ercole server user" ercole

%prep
rm -rf %{_topdir}/BUILD/%{name}-%{version}
cp -rf %{_topdir}/SOURCES/%{name}-%{version} %{_topdir}/BUILD/%{name}-%{version}
cd %{_topdir}/BUILD/%{name}-%{version}
chown -R root.root .
chmod -Rf a+rX,u+w,g-w,o-w .
cp target/%{name}-%{version}.jar %{name}.jar
cp package/rhel7/%{name}.service %{name}.service

%install
cd %{_topdir}/BUILD/%{name}-%{version}
mkdir -p %{buildroot}/opt/%{name}/run
install -m 0755 %{name}.jar %{buildroot}/opt/%{name}/%{name}.jar
mkdir -p %{buildroot}%{_unitdir} %{buildroot}%{_presetdir}
install -m 0644 package/rhel7/ercole-server.service %{buildroot}%{_unitdir}/%{name}.service
install -m 0644 package/rhel7/60-ercole-server.preset %{buildroot}%{_presetdir}/60-%{name}.preset

%post
/usr/bin/systemctl preset %{name}.service >/dev/null 2>&1 ||:

%preun
/usr/bin/systemctl --no-reload disable %{name}.service >/dev/null 2>&1 || :
/usr/bin/systemctl stop %{name}.service >/dev/null 2>&1 ||:

%postun
/usr/bin/systemctl daemon-reload >/dev/null 2>&1 ||:

%files
%attr(-,ercole,-) /opt/ercole-server/run
%dir /opt/ercole-server
/opt/ercole-server/ercole-server.jar
%{_unitdir}/ercole-server.service
%{_presetdir}/60-ercole-server.preset

%changelog
* Mon Aug 2 2019 Andrea Laisa <alaisa@sorint.it>
- 
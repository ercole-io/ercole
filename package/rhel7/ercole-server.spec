Name:           ercole-server
Version:        ERCOLE_VERSION
Release:        1%{?dist}
Summary:        Ercole server	

License:        ASL 2.0
URL:            https://ercole.io            
Source0:        https://github.com/ercole-io/%{name}/archive/%{version}.tar.gz
Group:          Tools
Requires:       java-11-openjdk
BuildRequires:  systemd systemd-rpm-macros

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
mkdir -p %{buildroot}/opt/%{name}/run %{buildroot}/etc/systemd/system
install -m 0755 %{name}.jar %{buildroot}/opt/%{name}/%{name}.jar
install -m 0644 %{name}.service %{_unitdir}/%{name}.service

%post
%systemd_post ercole.service

%preun
%systemd_preun ercole.service

%postun
%systemd_postun_with_restart ercole.service

%files
%attr(-,ercole,-) /opt/ercole-server/run
%dir /opt/ercole-server
/opt/ercole-server/ercole-server.jar
/etc/systemd/system/ercole-server.service

%changelog
* Mon Aug 2 2019 Andrea Laisa <alaisa@sorint.it>
- 
 
Name:           ercole
Version:        %{_version}
Release:        1%{?dist}
Summary:        Ercole	

License:        GPLv3
URL:            https://ercole.io            
Source0:        https://github.com/ercole-io/%{name}/archive/%{version}.tar.gz
Group:          Daemons
Requires:       systemd createrepo
BuildRequires:  systemd

%global         debug_package %{nil}

%description
Ercole is the server component of the ercole project.

%global debug_package %{nil}

%pre
    getent passwd ercole >/dev/null || useradd -s /bin/bash -c "Ercole user" ercole

%prep
rm -rf %{_builddir}/%{name}-%{version}
cp -rf %{_sourcedir}/%{name}-%{version} %{_builddir}/%{name}-%{version}
cd %{_builddir}/%{name}-%{version}
ls

%install
cd %{_builddir}/%{name}-%{version}
mkdir -p %{buildroot}/usr/bin/ %{buildroot}/usr/share/ercole %{buildroot}%{_unitdir} %{buildroot}%{_presetdir} %{buildroot}/var/lib/ercole/distributed_files/shared
install -m 0755 ercole %{buildroot}/usr/bin/ercole
install -m 0644 resources/initial_oracle_licenses_list.txt %{buildroot}/usr/share/ercole
install -m 0644 -d resources/templates %{buildroot}/usr/share/ercole

install -m 0644 package/config.toml %{buildroot}/usr/share/ercole/config.toml
install -m 0644 package/systemd/*.service %{buildroot}%{_unitdir}/
install -m 0644 package/systemd/60-ercole.preset %{buildroot}%{_presetdir}/60-%{name}.preset

%post
/usr/bin/systemctl preset %{name}.service >/dev/null 2>&1 ||:

%preun
/usr/bin/systemctl --no-reload disable %{name}.service >/dev/null 2>&1 || :
/usr/bin/systemctl stop %{name}.service >/dev/null 2>&1 ||:

%postun
/usr/bin/systemctl daemon-reload >/dev/null 2>&1 ||:

%files
%dir /var/lib/ercole
%dir /var/lib/ercole/distributed_files
/usr/bin/ercole
%{_presetdir}/60-ercole.preset
%{_unitdir}/ercole-alertservice.service
%{_unitdir}/ercole-apiservice.service
%{_unitdir}/ercole-dataservice.service
%{_unitdir}/ercole-reposervice.service
%{_unitdir}/ercole.service
/usr/share/ercole/config.toml
/usr/share/ercole/initial_oracle_licenses_list.txt

%changelog

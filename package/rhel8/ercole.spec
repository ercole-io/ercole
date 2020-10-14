 
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
echo "Creating ercole user..."
getent group ercole || groupadd -r ercole
getent passwd ercole >/dev/null || useradd -r -g ercole -s /bin/bash -c "Ercole user" ercole

%prep
rm -rf %{_builddir}/%{name}-%{version}
cp -rf %{_sourcedir}/%{name}-%{version} %{_builddir}/%{name}-%{version}
cd %{_builddir}/%{name}-%{version}
ls

%install
cd %{_builddir}/%{name}-%{version}
mkdir -p %{buildroot}/usr/bin/ %{buildroot}/usr/share/ercole/{examples,templates} %{buildroot}/usr/share/ercole/technologies/{Microsoft,Oracle,HP,IBM,RedHat,MariaDBFoundation,PostgreSQL,Unknown,VMWare} %{buildroot}%{_unitdir} %{buildroot}%{_presetdir} %{buildroot}/var/lib/ercole/distributed_files
install -m 0755 ercole %{buildroot}/usr/bin/ercole
install -m 0755 package/ercole-setup %{buildroot}/usr/bin/ercole-setup
install -m 0644 package/config.toml %{buildroot}/usr/share/ercole/config.toml
install -m 0644 resources/templates/* %{buildroot}/usr/share/ercole/templates/
install -m 0644 resources/technologies/list.json %{buildroot}/usr/share/ercole/technologies/list.json
install -m 0644 resources/technologies/Oracle/* %{buildroot}/usr/share/ercole/technologies/Oracle/
install -m 0644 resources/technologies/Microsoft/* %{buildroot}/usr/share/ercole/technologies/Microsoft/
install -m 0644 resources/technologies/HP/* %{buildroot}/usr/share/ercole/technologies/HP/
install -m 0644 resources/technologies/IBM/* %{buildroot}/usr/share/ercole/technologies/IBM/
install -m 0644 resources/technologies/RedHat/* %{buildroot}/usr/share/ercole/technologies/RedHat/
install -m 0644 resources/technologies/MariaDBFoundation/* %{buildroot}/usr/share/ercole/technologies/MariaDBFoundation/
install -m 0644 resources/technologies/PostgreSQL/* %{buildroot}/usr/share/ercole/technologies/PostgreSQL/
install -m 0644 resources/technologies/VMWare/* %{buildroot}/usr/share/ercole/technologies/VMWare/
install -m 0644 resources/technologies/Unknown/* %{buildroot}/usr/share/ercole/technologies/Unknown/
install -m 0644 package/systemd/*.service %{buildroot}%{_unitdir}/
install -m 0644 package/systemd/60-ercole.preset %{buildroot}%{_presetdir}/60-%{name}.preset
install -m 0644 distributed_files/ping.txt %{buildroot}/var/lib/ercole/distributed_files/ping.txt
install -m 0644 distributed_files/shared/*.repo %{buildroot}/usr/share/ercole/examples/

%post
echo "Running systemctl commands"
/usr/bin/systemctl daemon-reload >/dev/null 2>&1 ||:
/usr/bin/systemctl preset %{name}.service >/dev/null 2>&1 ||:
/usr/bin/systemctl preset %{name}-alertservice.service >/dev/null 2>&1 ||:
/usr/bin/systemctl preset %{name}-apiservice.service >/dev/null 2>&1 ||:
/usr/bin/systemctl preset %{name}-chartservice.service >/dev/null 2>&1 ||:
/usr/bin/systemctl preset %{name}-reposervice.service >/dev/null 2>&1 ||:
/usr/bin/systemctl preset %{name}-dataservice.service >/dev/null 2>&1 ||:
/usr/bin/systemctl is-active --quiet %{name}-alertservice.service && /usr/bin/systemctl restart %{name}-alertservice.service
/usr/bin/systemctl is-active --quiet %{name}-apiservice.service && /usr/bin/systemctl restart %{name}-apiservice.service
/usr/bin/systemctl is-active --quiet %{name}-chartservice.service && /usr/bin/systemctl restart %{name}-chartservice.service
/usr/bin/systemctl is-active --quiet %{name}-reposervice.service && /usr/bin/systemctl restart %{name}-reposervice.service
/usr/bin/systemctl is-active --quiet %{name}-dataservice.service && /usr/bin/systemctl restart %{name}-dataservice.service
echo "Running NOINTERACTIVE=1 /usr/bin/ercole-setup"
/bin/sh -c "NOINTERACTIVE=1 /usr/bin/ercole-setup"
echo "Running ercole completion bash"
ercole completion bash > /usr/share/bash-completion/completions/ercole

%preun
/usr/bin/systemctl --no-reload disable %{name}.service >/dev/null 2>&1 || :
/usr/bin/systemctl stop %{name}.service >/dev/null 2>&1 ||:

%postun
/usr/bin/systemctl daemon-reload >/dev/null 2>&1 ||:

%files
%dir /var/lib/ercole/distributed_files
/usr/bin/ercole
/usr/bin/ercole-setup
%{_presetdir}/60-ercole.preset
%{_unitdir}/ercole-alertservice.service
%{_unitdir}/ercole-apiservice.service
%{_unitdir}/ercole-chartservice.service
%{_unitdir}/ercole-dataservice.service
%{_unitdir}/ercole-reposervice.service
%{_unitdir}/ercole.service
/usr/share/ercole/config.toml
/usr/share/ercole/technologies/list.json
/usr/share/ercole/technologies/Oracle/Database.png
/usr/share/ercole/technologies/Oracle/Solaris.png
/usr/share/ercole/technologies/Oracle/MySQL.png
/usr/share/ercole/technologies/Oracle/VM.png
/usr/share/ercole/technologies/Oracle/Exadata.png
/usr/share/ercole/technologies/Microsoft/SQLServer.png
/usr/share/ercole/technologies/Microsoft/WindowsServer2008.png
/usr/share/ercole/technologies/Microsoft/WindowsServer2012.png
/usr/share/ercole/technologies/Microsoft/WindowsServer2016.png
/usr/share/ercole/technologies/Microsoft/WindowsServer2019.png
/usr/share/ercole/technologies/Microsoft/SQLServer.png
/usr/share/ercole/technologies/HP/HPUX.png
/usr/share/ercole/technologies/IBM/AIX.png
/usr/share/ercole/technologies/RedHat/EnterpriseLinux5.png
/usr/share/ercole/technologies/RedHat/EnterpriseLinux6.png
/usr/share/ercole/technologies/RedHat/EnterpriseLinux7.png
/usr/share/ercole/technologies/RedHat/EnterpriseLinux8.png
/usr/share/ercole/technologies/MariaDBFoundation/MariaDB.png
/usr/share/ercole/technologies/PostgreSQL/PostgreSQL.png
/usr/share/ercole/technologies/Unknown/Unknown.png
/usr/share/ercole/technologies/VMWare/VMWare.png
/usr/share/ercole/templates/template_addm.xlsx
/usr/share/ercole/templates/template_clusters.xlsx
/usr/share/ercole/templates/template_databases.xlsx
/usr/share/ercole/templates/template_hosts.xlsx
/usr/share/ercole/templates/template_lms.xlsm
/usr/share/ercole/templates/template_patch_advisor.xlsx
/usr/share/ercole/templates/template_segment_advisor.xlsx
/usr/share/ercole/templates/template_alerts.xlsx
/usr/share/ercole/examples/ercole-rhel5-x86_64.repo
/usr/share/ercole/examples/ercole-rhel6-x86_64.repo
/usr/share/ercole/examples/ercole-rhel7-x86_64.repo
/usr/share/ercole/examples/ercole-rhel8-x86_64.repo
/usr/share/ercole/examples/ercole-rhel7-noarch.repo
/usr/share/ercole/examples/ercole-rhel8-noarch.repo
/var/lib/ercole/distributed_files/ping.txt

%changelog

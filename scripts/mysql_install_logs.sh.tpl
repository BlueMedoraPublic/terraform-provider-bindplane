#!/bin/bash

set -e

sudo su

apt-get update

apt-get -y install wget
export COMPANY_ID=${company_id}
export SECRET_KEY=${secret_key}
export TEMPLATE_ID=${template_id}
wget https://bindplane-logs-downloads.s3.amazonaws.com/agent/unix_install.sh
chmod +x unix_install.sh
./unix_install.sh -a y

echo "mysql-server-5.7 mysql-server/root_password password ${mysql_pass}" | sudo debconf-set-selections
echo "mysql-server-5.7 mysql-server/root_password_again password ${mysql_pass}" | sudo debconf-set-selections
apt-get -y install mysql-server-5.7
sed -i 's/bind-address/#bind-address/g' /etc/mysql/mysql.conf.d/mysqld.cnf
echo "general_log = on" >> /etc/mysql/mysql.conf.d/mysqld.cnf
echo "general_log_file=/var/log/mysql/general.log" >> /etc/mysql/mysql.conf.d/mysqld.cnf
echo "log_error=/var/log/mysql/mysqld.log" >> /etc/mysql/mysql.conf.d/mysqld.cnf
echo "slow_query_log = on" >> /etc/mysql/mysql.conf.d/mysqld.cnf
echo "slow_query_log_file=/var/log/mysql/slow.log" >> /etc/mysql/mysql.conf.d/mysqld.cnf
mkdir -p /var/log/mysql
chown -R mysql:mysql /var/log/mysql
systemctl enable mysql
systemctl restart mysql
cat << EOF > /tmp/setup.sql
CREATE DATABASE ${database};
GRANT ALL PRIVILEGES ON *.* TO '${mysql_user}' IDENTIFIED BY '${mysql_pass}';
FLUSH PRIVILEGES;
EOF
mysql -u root -p${mysql_pass} < /tmp/setup.sql >> /tmp/setup.log 2>&1

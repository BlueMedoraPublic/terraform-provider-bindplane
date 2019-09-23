#!/bin/bash

sudo apt-get update
sudo ufw allow 3306
echo "mysql-server-5.7 mysql-server/root_password password ${mysql_pass}" | sudo debconf-set-selections
echo "mysql-server-5.7 mysql-server/root_password_again password ${mysql_pass}" | sudo debconf-set-selections
sudo apt-get -y install mysql-server-5.7 >> /tmp/setup.log 2>&1
sudo sed -i 's/bind-address/#bind-address/g' /etc/mysql/mysql.conf.d/mysqld.cnf
sudo systemctl enable mysql >> /tmp/setup.log 2>&1
sudo systemctl restart mysql >> /tmp/setup.log 2>&1
cat << EOF > /tmp/setup.sql
CREATE DATABASE ${database};
GRANT ALL PRIVILEGES ON *.* TO '${mysql_user}' IDENTIFIED BY '${mysql_pass}';
FLUSH PRIVILEGES;
EOF
mysql -u root -p${mysql_pass} < /tmp/setup.sql >> /tmp/setup.log 2>&1

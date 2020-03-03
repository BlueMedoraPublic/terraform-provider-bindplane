#!/bin/bash

sudo apt-get update
sudo ufw allow 3306
echo "mysql-server-5.7 mysql-server/root_password password ${mysql_pass}" | sudo debconf-set-selections
echo "mysql-server-5.7 mysql-server/root_password_again password ${mysql_pass}" | sudo debconf-set-selections
sudo apt-get -y install wget mysql-server-5.7 >> /tmp/setup.log 2>&1
sudo sed -i 's/bind-address/#bind-address/g' /etc/mysql/mysql.conf.d/mysqld.cnf
sudo systemctl enable mysql >> /tmp/setup.log 2>&1
sudo systemctl restart mysql >> /tmp/setup.log 2>&1
cat << EOF > /tmp/setup.sql
CREATE DATABASE ${database};
GRANT ALL PRIVILEGES ON *.* TO '${mysql_user}' IDENTIFIED BY '${mysql_pass}';
FLUSH PRIVILEGES;
EOF
mysql -u root -p${mysql_pass} < /tmp/setup.sql >> /tmp/setup.log 2>&1

export COMPANY_ID=${company_id}
export SECRET_KEY=${secret_key}
export TEMPLATE_ID=${template_id}

wget https://bindplane-logs-downloads.s3.amazonaws.com/agent/unix_install.sh
chmod +x unix_install.sh
sudo ./unix_install.sh -a y

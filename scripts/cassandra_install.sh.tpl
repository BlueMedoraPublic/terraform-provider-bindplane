yum install -y wget

# install java runtime
wget \
    --no-cookies --no-check-certificate \
    --header "Cookie:oraclelicense=accept-securebackup-cookie" \
    "http://download.oracle.com/otn-pub/java/jdk/8u131-b11/d54c1d3a095b4ff2b6607d096fa80163/jdk-8u131-linux-x64.rpm"

yum -y localinstall jdk-8u131-linux-x64.rpm


# set java home for root user
export JAVA_HOME=/usr/java/jdk1.8.0_131/ && echo 'export JAVA_HOME=/usr/java/jdk1.8.0_131' >> /root/.bashrc
export JRE_HOME=/usr/java/jdk1.8.0_131/jre && echo 'export JRE_HOME=/usr/java/jdk1.8.0_131/jre' >> /root/.bashrc


# setup repo
cat <<EOF > /etc/yum.repos.d/cassandra.repo
[cassandra]
name=Apache Cassandra
baseurl=https://www.apache.org/dist/cassandra/redhat/311x/
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://www.apache.org/dist/cassandra/KEYS
EOF


# install cassandra
yum -y install cassandra


# setup jmx
sed -i "s/LOCAL_JMX=yes/LOCAL_JMX=no/g" /etc/cassandra/conf/cassandra-env.sh
cat <<EOF > /etc/cassandra/jmxremote.password
cassandra cassandra
EOF
cat <<EOF > /etc/cassandra/jmxremote.access
cassandra readwrite
EOF
chmod 400 /etc/cassandra/jmxremote.password && chown cassandra:cassandra /etc/cassandra/jmxremote.password
chmod 400 /etc/cassandra/jmxremote.access && chown cassandra:cassandra /etc/cassandra/jmxremote.access


# service restart after updating config
sudo systemctl enable cassandra
sudo systemctl restart cassandra

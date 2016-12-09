FROM mariadb:5.5
ADD create.sql /docker-entrypoint-initdb.d/
ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_DATABASE=platform
ENV MYSQL_USER=chat
ENV MYSQL_PASSWORD=chatpwd

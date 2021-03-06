FROM ubuntu:18.04

RUN apt-get -y update
ENV USERNAME toringol
ENV PASSWORD toringol

#
# Установка postgresql
#
ENV PGVER 10
RUN apt-get install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER $USERNAME WITH SUPERUSER PASSWORD '$PASSWORD';" &&\
    createdb -O $USERNAME forumdb &&\
    psql --command "CREATE EXTENSION  IF NOT EXISTS citext;" &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf


# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432
USER root
# Add VOLUMEs to allow backup of config, logs and databases
#VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Установка golang
ENV GOVER 1.10
RUN apt-get install -y golang-$GOVER
RUN apt-get install -y git

# Выставляем переменную окружения для сборки проекта
ENV GOROOT /usr/lib/go-$GOVER
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

# Копируем исходный код в Docker-контейнер
WORKDIR $GOPATH/src/github.com/Toringol/forum
COPY . $GOPATH/src/github.com/Toringol/forum

RUN go install .
EXPOSE 5000
#RUN /etc/init.d/postgresql start &&\
#    psql -U $USERNAME -d codeloft -a -f resources/initdb.sql &&\
#    /etc/init.d/postgresql stop
CMD service postgresql start && forum



FROM postgres:11.5-alpine

USER postgres

RUN chmod 0700 /var/lib/postgresql/data &&\
    initdb /var/lib/postgresql/data &&\
    echo "host all  all    0.0.0.0/0  md5" >> /var/lib/postgresql/data/pg_hba.conf &&\
    echo "listen_addresses='*'" >> /var/lib/postgresql/data/postgresql.conf &&\
    pg_ctl start

EXPOSE 5432

# ENV GOSU_VERSION 1.11
# RUN set -eux; \
# 	\
# 	apk add --no-cache --virtual .gosu-deps \
# 		ca-certificates \
# 		dpkg \
# 		gnupg \
# 	; \
# 	\
# 	dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')"; \
# 	wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch"; \
# 	wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch.asc"; \
# 	\
#   # verify the signature
# 	export GNUPGHOME="$(mktemp -d)"; \
# 	gpg --batch --keyserver hkps://keys.openpgp.org --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4; \
# 	gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu; \
# 	command -v gpgconf && gpgconf --kill all || :; \
# 	rm -rf "$GNUPGHOME" /usr/local/bin/gosu.asc; \
# 	\
#   # clean up fetch dependencies
# 	apk del --no-network .gosu-deps; \
# 	\
# 	chmod +x /usr/local/bin/gosu; \
#   # verify that the binary works
# 	gosu --version; \
# 	gosu nobody true

# RUN chmod +x docker-entrypoint.sh
USER root
RUN apk update && apk add --no-cache bash
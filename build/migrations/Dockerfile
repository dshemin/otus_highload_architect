FROM gomicro/goose

ADD /migrations/*.sql /migrations

ADD ./build/migrations/entrypoint.sh /migrations/entrypoint.sh

ENTRYPOINT ["bash", "/migrations/entrypoint.sh"]
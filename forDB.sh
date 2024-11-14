#!/bin/bash

psql -U admin <<-EOSQL
    CREATE DATABASE "ozonTest";
    CREATE DATABASE "ozonTest-test";
EOSQL
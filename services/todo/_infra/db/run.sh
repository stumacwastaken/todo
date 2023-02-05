#!/bin/sh

# Get number of "up" files to find proper number
upCount=$(find migrations -type f | grep ".up.sql" | wc -l)

printf "found ${upCount} up migration files \n"
printf "applying ${upCount} migrations to host ${HOST} as ${USER} with ${PASSWORD} \n"

migrate -database "mysql://${USER}:${PASSWORD}@tcp(${HOST})/${DB}" -path ./migrations up 
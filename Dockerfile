from ubuntu
maintainer Florian Fink <finkf@cis.lmu.de>
copy lmdcache /app/lmdcache
cmd /app/lmdcache -host '0.0.0.0:8080' -lmd 'http://lmd:8181'

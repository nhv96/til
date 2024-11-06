# Clean up Docker unused resources
Sometime after using Docker I notice my disk is taking too much storage.

I ran command to check what is taking all the space:
```
$ sudo du -h /var | sort -h
...

1.4G	/var/cache
1.9G	/var/lib/snapd/cache
2.3G	/var/lib/docker/volumes
3.1G	/var/lib/snapd/snaps
4.0G	/var/log/journal
4.0G	/var/log/journal/29a2a98d1f5c4275b8474be92a9c1751
4.1G	/var/log
5.8G	/var/lib/snapd
29G	/var/lib/docker/overlay2
31G	/var/lib/docker
37G	/var/lib
42G	/var
```

Turn out docker `overlay2` is taking so much space, so I use this command to clean up unused resources.
```
docker system prune
...

6d46882219ug8hfuj8jpkeff3
iqc16xm2p8d7ppyvami96ny6v

Total reclaimed space: 21.52GB
```
# Summary
This code will cause kernel panic / totally lockup nfs4 server on newest FC 29. It'll also cause NFS Server to lock up / stop responding on older FC versions. Now why is it called drbd_kill? It seems that additional layer
which drbd is, is causing the code to kill the server faster. But it turned out it only needs nfs4 share to kill nfs4-server completely. NFS Client seems to be working ok in all cases. So this is purely NFS Server related.

Normally it takes about 20-30s, when you add drbd it'll take around 5-15s. For testing with drbd - you need to setup drbd volume, mount it, then share it via nfs4 (client and server can be on the same machine). Then you need to run the app on the NFS share.

It'll kill both NFS 4.1 and NFS 4.2. NFS 3 seems not affected.

## Build

<pre>
Fedora release 29 (Twenty Nine)
- 5.0.7-200.fc29.x86_64 #1 SMP (Fedora 29, also other sub-versions)
- 5.0.6-200.fc29.x86_64 #1 SMP (Fedora 29 Server Edition)

cat /proc/drbd
version: 8.4.11 (api:1/proto:86-101)
srcversion: C27D50EE6C67ED861348AA6
</pre>

Compilation instruction:
<pre>
1. export GOPATH=/nfsshare/
2. cd /nfsshare
3. mkdir src
4. put the sources into src dir
5. cd src/multiaccess
6. go build *.go
7. run ./multiaccess couple of times (usually takes 2-3 runs to kill the server)
</pre>

## How it works
This code will create 30 threads which will randomly read/write/seek/lock/rename/unlink 30 random files for 5 seconds. Probably needs to be run on SSD. Now these random operations are all working fine beside flock with LOCK_EX, even LOCK_SH is working correctly and if you replace exclusive locks with shared locks this app shouldn't be able to kill the server. LOCK_EX is the thing which will make the server to crash under high load.

Crash log is attached.

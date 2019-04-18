# Summary
This code will cause kernel panic / totally lockup nfs4 server on newest FC 29. Now why is it called drbd_kill? It seems that additional layer
which drbd is, is causing the code to kill the server faster. But it turned out it only needs nfs4 share to kill nfs4-server completely. Client seems to be working ok.

Normally it takes about 20-30s, when you add drbd it'll take around 5-15s. For testing with drbd - you need to setup drbd volume, mount it, then share it via nfs4 (client and server can be on the same machine). Then you need to run the app on the NFS share.

## Build

<pre>
Fedora release 29 (Twenty Nine)
Kernel:  5.0.7-200.fc29.x86_64 #1 SMP

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
This code will create 30 threads which will randomly read/write/lock/rename 30 random files for 5 seconds. Probably needs to be run on SSD
to properly kill the server (wasn't tested on HDD and it might be too slow)

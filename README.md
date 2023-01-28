# Cgstat

Cgstat is a tool for displaying stats about cgroups. 

## Getting started
Install `cgstat` on your local host using: 
```bash
make install
```
If you're a developer and want to test out changes that you're making,
then build `cgstat` with the following command:

```bash
make build
```

Before submitting a pull request, be sure to run

```bash
make release
```

## Usage

### Listing cgroups on a host
~~~~
$ cgstat list
~~~~

### Listing cgroups by prefix
If you're only interested in a subset of cgroups, then you can filter them by
prefix using the `--prefix` parameter.
~~~~
$ cgstat list --prefix=
~~~~

### Viewing cgroups on a host
```
# View stats on a specific cgroup
$ cgstat view --name=/system.slice/sshd.service 
Name                        UserCPU  KernelCPU  CurrentUsage      MaxUsage          UsageLimit  RSS        Cache      Dirty  WriteBack  UnderOom  OomKill  
/system.slice/sshd.service  53.67%   42.19%     27.7 MiB (0.00%)  31.3 MiB (0.00%)  8.0 EiB     908.0 KiB  132.0 KiB  0 B    0 B        0         0

# View verbose information about a given cgroup
cgstat --name=/system.slice/sshd.service --verbose

# View stats about a set of cgroups by prefix  
$ cgstat --prefix=/system.slice

# Follow updates in real time
$ cgstat --prefix=/system.slice --follow 
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[Apache 2.0](https://choosealicense.com/licenses/apache-2.0/)
# Cgstat

Cgstat is a tool for displaying stats about cgroups. 

## Installation

```bash
make build
```

## Usage

```python
# View stats on a specific cgroup
$ cgstat --name=/system.slice/sshd.service 
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
`pct` calculates the distribution of lines in text. It is similar to `sort | uniq -c | sort -n -r`, except that it prints percentages as well as counts.

Sample output:

```bash
$ cat std.qc | grep "clearfat" | pct | head
 33.05%    1204    clearfat q=2 c=0
 18.80%    685     clearfat q=4 c=0
 17.76%    647     clearfat q=3 c=0
  4.47%    163     clearfat q=6 c=0
  3.90%    142     clearfat q=8 c=0
  3.84%    140     clearfat q=10 c=0
  2.72%    99      clearfat q=5 c=0
  1.95%    71      clearfat q=13 c=0
  1.73%    63      clearfat q=4 c=4
  1.40%    51      clearfat q=7 c=0
```

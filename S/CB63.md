
# String functions

## Documentations

`go doc github.com/kokizzu/gotro/S`

## S.EncodeCB63, S.DecodeCB63 
convert integer from/to url-safe string

```
        112747872621 to      0d-I-4h back to         112747872621
 8361608650129439717 to  6F9OgExn6j_ back to  8361608650129439717
             1036326 to         2x-a back to              1036326
 1990494840112300345 to  0iUe7AF833t back to  1990494840112300345
 5237504367641861460 to  3XjLbZDqOpJ back to  5237504367641861460
       4354798931041 to      zMiEZlW back to        4354798931041
 7937896837161115613 to  5sd3lrPfpzS back to  7937896837161115613
       9543511901307 to     19s4ckGv back to        9543511901307
      12662292403333 to     1sFfRCH4 back to       12662292403333
             4030758 to         EN3a back to              4030758
 1488531474229776968 to  0HcJroyS487 back to  1488531474229776968
                 127 to  ---------0z back to                  127
               32767 to  --------6zz back to                32767
          2147483647 to  -----0zzzzz back to           2147483647
 9223372036854775807 to  6zzzzzzzzzz back to  9223372036854775807
                   0 to  ----------- back to                    0
                   1 to  ----------0 back to                    1
                   9 to  ----------8 back to                    9
                  10 to  ----------9 back to                   10
                  99 to  ---------0Y back to                   99
                 100 to  ---------0Z back to                  100
                 999 to  ---------Eb back to                  999
                1000 to  ---------Ec back to                 1000
                9999 to  --------1RE back to                 9999
               10000 to  --------1RF back to                10000
               99999 to  --------NPU back to                99999
              100000 to  --------NPV back to               100000
              999999 to  -------2o7z back to               999999
             1000000 to  -------2o8- back to              1000000
             9999999 to  -------a8Oz back to              9999999
            10000000 to  -------a8P- back to             10000000
            99999999 to  ------4xT2z back to             99999999
           100000000 to  ------4xT3- back to            100000000
           999999999 to  ------vagbz back to            999999999
          1000000000 to  ------vagc- back to           1000000000
          9999999999 to  -----8J1yEz back to           9999999999
         10000000000 to  -----8J1yF- back to          10000000000
         99999999999 to  ----0S7SiUz back to          99999999999
        100000000000 to  ----0S7SiV- back to         100000000000
        999999999999 to  ----DYJdFzz back to         999999999999
       1000000000000 to  ----DYJdG-- back to        1000000000000
       9999999999999 to  ---1GWDRdzz back to        9999999999999
      10000000000000 to  ---1GWDRe-- back to       10000000000000
      99999999999999 to  ---LjBFTYzz back to       99999999999999
     100000000000000 to  ---LjBFTZ-- back to      100000000000000
     999999999999999 to  --2YMuZlbzz back to      999999999999999
    1000000000000000 to  --2YMuZlc-- back to     1000000000000000
    9999999999999999 to  --YWj8jkEzz back to     9999999999999999
   10000000000000000 to  --YWj8jkF-- back to    10000000000000000
   99999999999999999 to  -4YGMWSXUzz back to    99999999999999999
  100000000000000000 to  -4YGMWSXV-- back to   100000000000000000
  999999999999999999 to  -rVhfDbNzzz back to   999999999999999999
 1000000000000000000 to  -rVhfDbO--- back to  1000000000000000000
PASS
ok      github.com/kokizzu/gotro/S      0.803s
```
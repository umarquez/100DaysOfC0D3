# Diagramas

## 2 TEA rounds

```
 
 ---------------------------------------------
|   ------                                    |
|  | << n | = left shift         ( + ) = XOR  |
|   ------                                    |
|   ------                        ---         |
|  | >> n | = right shift        | + | = ADD  |
|   ------                        ---         |
 -------------------[ TEA ]-------------------


  |                                         |
  v                                         v
  |               K[0]                      |
  |                 |                       |
  |                 v                       |
  |                ---     ------           |
  |         ------| + |<--| << 4 |<---      |
  |        |       ---     ------     |     |
  |        |                          |     |
  |        |         --- Delta(i)     |     |
  |        |        |                 |     |
  v        |        v                 |     |
 ---       v       ---                |     |
| + |<---( + )<---| + |<--------------+-----|
 ---       ^       ---                |     |
  |        |                          |     |
  |        |       ---     ------     |     |
  |         ------| + |<--| >> 5 |<---      |
  |                ---     ------           |
  |                 ^                       |
  v                 |                       v
  |               K[1]                      |
   \                                       / 
    \             -->B     A<--           /  
     -------->---------  --------<--------
                       \/
                  <--A /\  B-->
     --------<---------  -------->--------
    /                                     \
   /                                       \
  |                                         |
  v               K[2]                      v
  |                 |                       |
  |                 v                       |
  |                ---     ------           |
  |         ------| + |<--| << 4 |<---      |
  |        |       ---     ------     |     |
  |        |                          |     |
  |        |         --- Delta(i)     |     |
  |        |        |                 |     |
  v        |        v                 |     |
 ---       v       ---                |     |
| + |<---( + )<---| + |<--------------+-----|
 ---       ^       ---                |     |
  |        |                          |     |
  |        |       ---     ------     |     |
  |         ------| + |<--| >> 5 |<---      |
  |                ---     ------           |
  |                 ^                       |
  v                 |                       v
  |               K[3]                      |
   \                                       / 
    \             -->B     A<--           /  
     -------->---------  --------<--------
                       \/
                  <--A /\  B-->
     --------<---------  -------->--------
    /                                     \
   /                                       \
  |                                         |
  v                                         v
  |                                         |
  
```
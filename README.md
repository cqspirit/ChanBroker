# ChanBroker

# �����ܲ���ʾ�� 

## д�ļ� 
```
root@ubuntu:/share/go/src/github.com/myself659/ChanBroker# time  go  run example/profile.go   -n 2000  -t  60    > profile.txt  
real    1m11.178s
user    0m45.708s
sys     2m13.924s
root@ubuntu:/share/go/src/github.com/myself659/ChanBroker# tail     profile.txt  
0xc42005a1e0 event: 21827
0xc42005a1e0 event: 21828
0xc42005a1e0 event: 21829
0xc42005a1e0 event: 21830
0xc42005a1e0 event: 21831
0xc42005a1e0 event: 21832
0xc42005a1e0 event: 21833
0xc42005a1e0 event: 21834
0xc42005a1e0 has recv: 21835
total: 43670000

```

��real time�������ܽ�����£�

ÿ�뷢���������Ϣ������ 43670000/72= 606527  
ÿ����Է�����Ϣ������21835/60 = 363 

### ��д�ļ�

```
root@ubuntu:/share/go/src/github.com/myself659/ChanBroker# time  go  run example/profile.go   -n 2000  -t  60  

������
0xc42007c120 has recv: 159444
total: 318888000

real    1m7.561s
user    2m21.824s
sys     0m20.936s
```

��user time�������ܽ�����£�

ÿ�뷢���������Ϣ������ 318888000/142= 2445690 
ÿ����Է�����Ϣ������159444/60 = 2657   
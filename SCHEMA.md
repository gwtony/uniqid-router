## Data between Agent and Router
```
raw data:
magic: 0x77 
len: uint16_t  //len of data
data: uniqid_data

uniquid_data {
	uint8_t    uid[32];
	uint8_t    puid[32];
	uint8_t    pip[4];
	uint16_t  pport;
	uint8_t    lip[4];
	uint16_t  lport;
	uint16_t  dlen;
	uint8_t    data[0];
}
```

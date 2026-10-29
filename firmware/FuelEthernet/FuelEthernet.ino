// Cheesy Arena F2 protocol v1. Requires Ethernet 2.0.2 and Crypto 0.4.0.
// No heap Strings, JSON or independent match clock. ISR queue size is bounded.
#include <SPI.h>
#include <Ethernet.h>
#include <SHA256.h>
#include <EEPROM.h>
#include <util/atomic.h>
#include "DeviceConfig.h"

const uint8_t FRAME_SIZE=88, QUEUE_SIZE=16;
const uint8_t HELLO=1, CHALLENGE=2, AUTH=3, STATE=4, EVENT=5, ACK=6, HEARTBEAT=7;
const uint32_t FIRMWARE_VERSION=0x010000;
EthernetClient client;
SHA256 hash;
uint8_t rx[FRAME_SIZE], tx[FRAME_SIZE], nonce[16], attempt[16];
uint8_t received=0;
uint32_t sendSequence=0, receiveSequence=0, bootId=0;
uint32_t lastRx=0, lastPing=0, lastConnect=0, sentAt=0;
bool authenticated=false, stateSynced=false, waitingAck=false;
uint32_t waitingEvent=0;
// Latest authenticated server state for optional, board-specific indicator code.
uint32_t matchFlags=0, matchId=0;
uint8_t matchPhase=0, scoringShift=8;
struct SensorEvent { uint32_t number; uint32_t lease; };
volatile SensorEvent queue[QUEUE_SIZE];
volatile uint8_t head=0, tail=0;
volatile uint32_t detected=0, dropped=0, activeLease=0, lastEdge[2]={0,0};
volatile bool captureReady=false;

uint32_t get32(const uint8_t *p) { return ((uint32_t)p[0]<<24)|((uint32_t)p[1]<<16)|((uint32_t)p[2]<<8)|p[3]; }
void put32(uint8_t *p,uint32_t v) { p[0]=v>>24; p[1]=v>>16; p[2]=v>>8; p[3]=v; }

void sensorEdge(uint8_t channel) {
  uint32_t now=micros();
  if (lastEdge[channel] && (uint32_t)(now-lastEdge[channel])<SENSOR_DEBOUNCE_US) return;
  lastEdge[channel]=now;
  if (detected==0xFFFFFFFFUL) { captureReady=false; return; }
  uint32_t number=++detected;
  uint8_t next=(head+1)%QUEUE_SIZE;
  if (!captureReady || next==tail) { ++dropped; return; }
  queue[head].number=number; queue[head].lease=activeLease; head=next;
}
void sensor0() { sensorEdge(0); }
void sensor1() { sensorEdge(1); }

// Wear-levelled boot IDs. A complete complement pair is committed before use;
// a brownout during the next slot leaves the previous boot record intact.
uint32_t nextBootId() {
  uint32_t highest=0; uint8_t nextSlot=0;
  for(uint8_t slot=0;slot<16;++slot) {
    uint32_t value,check; EEPROM.get(slot*8,value); EEPROM.get(slot*8+4,check);
    if(value>0 && value<0xFFFFFFFEUL && check==~value && value>highest) { highest=value; nextSlot=(slot+1)%16; }
  }
  if(highest>=0xFFFFFFFDUL) return 0;
  uint32_t value=highest+1,invalid=0;
  EEPROM.put(nextSlot*8+4,invalid); EEPROM.put(nextSlot*8,value);
  uint32_t check=~value; EEPROM.put(nextSlot*8+4,check);
  return value;
}

void signFrame(uint8_t *frame) {
  hash.resetHMAC(DEVICE_KEY,sizeof(DEVICE_KEY)); hash.update(frame,72);
  hash.finalizeHMAC(DEVICE_KEY,sizeof(DEVICE_KEY),frame+72,16);
}
bool validFrame() {
  if(rx[0]!='F'||rx[1]!='2'||rx[2]!=1) return false;
  uint8_t tag[16], difference=0;
  hash.resetHMAC(DEVICE_KEY,sizeof(DEVICE_KEY)); hash.update(rx,72);
  hash.finalizeHMAC(DEVICE_KEY,sizeof(DEVICE_KEY),tag,16);
  for(uint8_t i=0;i<16;++i) difference|=tag[i]^rx[72+i];
  for(uint8_t i=0;i<16;++i) {
    uint8_t expected=i<strlen(DEVICE_ID)?DEVICE_ID[i]:0;
    difference|=rx[40+i]^expected;
  }
  return difference==0;
}

void disconnectDevice() {
  ATOMIC_BLOCK(ATOMIC_RESTORESTATE) {captureReady=false;}
  client.stop(); authenticated=false; stateSynced=false; received=0; waitingAck=false;
}

bool sendFrame(uint8_t type,uint32_t value=0,uint32_t lease=0,uint32_t flags=0) {
  memset(tx,0,sizeof(tx)); tx[0]='F';tx[1]='2';tx[2]=1;tx[3]=type;
  if(type!=HELLO) { if(sendSequence==0xFFFFFFFFUL) {disconnectDevice();return false;} ++sendSequence; }
  put32(tx+4,sendSequence); memcpy(tx+8,attempt,16); put32(tx+24,lease);
  uint32_t total; ATOMIC_BLOCK(ATOMIC_RESTORESTATE) {total=detected;}
  put32(tx+28,total);put32(tx+32,flags);put32(tx+36,value);
  memcpy(tx+40,DEVICE_ID,strlen(DEVICE_ID));memcpy(tx+56,nonce,16);signFrame(tx);
  if(client.write(tx,sizeof(tx))!=sizeof(tx)) {disconnectDevice();return false;}
  return true;
}

void handleFrame() {
  if(!validFrame() || receiveSequence==0xFFFFFFFFUL || get32(rx+4)!=receiveSequence+1) {disconnectDevice();return;}
  ++receiveSequence;
  if(rx[3]==CHALLENGE && !authenticated) {
    memcpy(nonce,rx+56,16);
    if(sendFrame(AUTH,bootId)) authenticated=true;
    return;
  }
  if(!authenticated || memcmp(nonce,rx+56,16)!=0) {disconnectDevice();return;}
  if(rx[3]==STATE) {
    matchFlags=get32(rx+32);matchId=get32(rx+28);
    matchPhase=(get32(rx+36)>>8)&255;scoringShift=get32(rx+36)&255;
    bool newContext=memcmp(attempt,rx+8,16)!=0;
    ATOMIC_BLOCK(ATOMIC_RESTORESTATE) {
      captureReady=false;
      if(newContext) {dropped+=(head+QUEUE_SIZE-tail)%QUEUE_SIZE;tail=head;}
      activeLease=get32(rx+24);
      captureReady=bootId!=0 && detected!=0xFFFFFFFFUL;
    }
    if(newContext) {memcpy(attempt,rx+8,16);waitingAck=false;}
    stateSynced=true;
    // Flags and phase are available here for optional indicator outputs.
    // Sensor capture never trusts local timing to decide whether a point counts.
  } else if(rx[3]==ACK && waitingAck && get32(rx+36)==waitingEvent) {
    if(get32(rx+32)==3) {disconnectDevice();return;} // Preserve queue on storage uncertainty.
    ATOMIC_BLOCK(ATOMIC_RESTORESTATE) {if(tail!=head && queue[tail].number==waitingEvent) tail=(tail+1)%QUEUE_SIZE;}
    waitingAck=false;
  }
}

void setup() {
  pinMode(SD_CS_PIN,OUTPUT);digitalWrite(SD_CS_PIN,HIGH);
  Ethernet.init(ETHERNET_CS_PIN);
  // The sensor input voltage must already be compatible with UNO logic levels.
  for(uint8_t i=0;i<SENSOR_CHANNELS && i<2;++i) pinMode(SENSOR_PINS[i],SENSOR_PULLUP?INPUT_PULLUP:INPUT);
  bootId=nextBootId();
  if(bootId==0) return;
  if(SENSOR_CHANNELS>0) attachInterrupt(digitalPinToInterrupt(SENSOR_PINS[0]),sensor0,SENSOR_ACTIVE_LOW?FALLING:RISING);
  if(SENSOR_CHANNELS>1) attachInterrupt(digitalPinToInterrupt(SENSOR_PINS[1]),sensor1,SENSOR_ACTIVE_LOW?FALLING:RISING);
  uint8_t mac[6];memcpy(mac,DEVICE_MAC,6);
  Ethernet.begin(mac,IPAddress(DEVICE_IP),IPAddress(GATEWAY_IP),IPAddress(GATEWAY_IP),IPAddress(SUBNET));
  Ethernet.setRetransmissionTimeout(200);Ethernet.setRetransmissionCount(2);
  client.setConnectionTimeout(300);
}

void loop() {
  if(bootId==0) return;
  uint32_t now=millis();
  if(!client.connected()) {
    ATOMIC_BLOCK(ATOMIC_RESTORESTATE) {captureReady=false;}
    if((uint32_t)(now-lastConnect)<1000) return;
    lastConnect=now;disconnectDevice();memset(nonce,0,16);sendSequence=receiveSequence=0;
    if(client.connect(IPAddress(SERVER_IP),SERVER_PORT)) {lastRx=now;sendFrame(HELLO,bootId,0,FIRMWARE_VERSION);}
    return;
  }
  // TCP is a byte stream: retain incomplete frames and process coalesced frames.
  while(client.available()) {
    int value=client.read();if(value<0) break;
    rx[received++]=(uint8_t)value;
    if(received==FRAME_SIZE) {received=0;handleFrame();lastRx=now;if(!client.connected()) return;}
  }
  if((uint32_t)(now-lastRx)>1000) {disconnectDevice();return;}
  if(waitingAck && (uint32_t)(now-sentAt)>1000) {disconnectDevice();return;}
  if(authenticated && (uint32_t)(now-lastPing)>=500) {
    uint32_t lost;ATOMIC_BLOCK(ATOMIC_RESTORESTATE){lost=dropped;}
    if(!sendFrame(HEARTBEAT,lost)) return;lastPing=now;
  }
  if(authenticated && stateSynced && !waitingAck) {
    SensorEvent event;bool available;
    ATOMIC_BLOCK(ATOMIC_RESTORESTATE) {available=tail!=head;if(available){event.number=queue[tail].number;event.lease=queue[tail].lease;}}
    if(available && sendFrame(EVENT,event.number,event.lease)) {waitingEvent=event.number;waitingAck=true;sentAt=now;}
  }
}

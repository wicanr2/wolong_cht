/* Original native C startup graph; re/133, spec/253. No CPU.Step or guest decoder. */
static int su_known(uint16_t t) {switch(t) {
case 0x0:
case 0x2:
case 0x4:
case 0x6:
case 0x9:
case 0xd:
case 0xf:
case 0x11:
case 0x13:
case 0x16:
case 0x19:
case 0x1e:
case 0x20:
case 0x22:
case 0x23:
case 0x24:
case 0x27:
case 0x29:
case 0x2b:
case 0x2c:
case 0x2d:
case 0x2f:
case 0x31:
case 0x34:
case 0x37:
case 0x3a:
case 0x3d:
case 0x40:
case 0x43:
case 0x46:
case 0x47:
case 0x49:
case 0x4b:
case 0x4e:
case 0x51:
case 0x54:
case 0x56:
case 0x57:
case 0x59:
case 0x5b:
case 0x5e:
case 0x61:
case 0x64:
case 0x67:
case 0x6a:
case 0x6b:
case 0x6c:
case 0x6f:
case 0x71:
case 0x74:
case 0x77:
case 0x7a:
case 0x7d:
case 0x7f:
case 0x82:
case 0x84:
case 0x86:
case 0x89:
case 0x8c:
case 0x8f:
case 0x90:
case 0x92:
case 0x95:
case 0x99:
case 0x9b:
case 0x9e:
case 0xa1:
case 0xa5:
case 0xa7:
case 0xa9:
case 0xac:
case 0xaf:
case 0xb2:
case 0xb5:
case 0xb9:
case 0xbb:
case 0xbd:
case 0xc0:
case 0xc3:
case 0xc6:
case 0xc9:
case 0xcd:
case 0xcf:
case 0xd2:
case 0xd5:
case 0xd9:
case 0xdb:
case 0xde:
case 0x155:
case 0x158:
case 0x15d:
case 0x160:
case 0x162:
case 0x167:
case 0x16a:
case 0x16d:
case 0x16f:
case 0x174:
case 0x177:
case 0x179:
case 0x17e:
case 0x181:
case 0x184:
case 0x189:
case 0x18c:
case 0x18f:
case 0x192:
case 0x197:
case 0x19a:
case 0x19d:
case 0x1a0:
case 0x1a5:
case 0x1a8:
case 0x1ab:
case 0x1ae:
case 0x1b3:
case 0x210:
case 0x211:
case 0x212:
case 0x213:
case 0x215:
case 0x217:
case 0x219:
case 0x21e:
case 0x220:
case 0x222:
case 0x224:
case 0x226:
case 0x22c:
case 0x232:
case 0x233:
case 0x234:
case 0x235:
case 0x236:
case 0x237:
case 0x239:
case 0x23b:
case 0x23d:
case 0x23f:
case 0x240:
case 0x33b:
case 0x33c:
case 0x33d:
case 0x33e:
case 0x340:
case 0x342:
case 0x345:
case 0x348:
case 0x34a:
case 0x34b:
case 0x34c:
case 0x34d:
case 0x34e:
case 0x34f:
case 0x352:
case 0x354:
case 0x355:
case 0x356:
case 0x357:
case 0x358:
case 0x35e:
case 0x364:
case 0x366:
case 0x36b:
case 0x36c:
case 0x612:
case 0x614:
case 0x616:
case 0x618:
case 0x61a:
case 0x61c:
case 0x61e:
case 0x1a6e:
case 0x1a70:
case 0x1a72:
case 0x1a74:
case 0x1a76:
case 0x1a79:
case 0x1a7b:
case 0x1a7e:
case 0x1a84:
case 0x1a87:
case 0x1a8a:
case 0x1a8d:
case 0x1a90:
case 0x1a93:
case 0x1a96:
case 0x1a99:
case 0x1a9c:
case 0x1a9e:
case 0x1aa1:
case 0x1aa4:
case 0x1aa7:
case 0x1aad:
case 0x1aae:
case 0x1ab1:
case 0x1ab4:
case 0x1ab7:
case 0x1aba:
case 0x1abc:
case 0x1abf:
case 0x1ac2:
case 0x1ac3:
case 0x1ac6:
case 0x1ac9:
case 0x1acc:
case 0x1acf:
case 0x1ad2:
case 0x1ad5:
case 0x1ad7:
case 0x1ad9:
case 0x1adc:
case 0x1ade:
case 0x1ae0:
case 0x1ae3:
case 0x1ae5:
case 0x1ae8:
case 0x1aea:
case 0x1aed:
case 0x1aef:
case 0x1af1:
case 0x1af3:
case 0x1af6:
case 0x1af8:
case 0x1afd:
case 0x1b00:
case 0x1b03:
case 0x1b05:
case 0x1b08:
case 0x1b0a:
case 0x1b0c:
case 0x1b0e:
case 0x1b10:
case 0x1b12:
case 0x1b17:
case 0x1b1c:
case 0x1b1e:
case 0x1b20:
case 0x1b25:
case 0x1b27:
case 0x1b29:
case 0x1b2c:
case 0x1b2f:
case 0x1b34:
case 0x1b39:
case 0x1b3c:
case 0x1b3e:
case 0x1b40:
case 0x1b42:
case 0x1b44:
case 0x1b48:
case 0x1b4c:
case 0x1b4f:
case 0x1b52:
case 0x1b55:
case 0x1b57:
case 0x1b59:
case 0x1be0:
case 0x1be2:
case 0x1be4:
case 0x1be6:
case 0x1be9:
case 0x1bec:
case 0x1bf1:
case 0x1bf6:
case 0x1bf9:
case 0x1bfd:
case 0x1c01:
case 0x1c04:
case 0x1c07:
case 0x1c0a:
case 0x1c0d:
case 0x1c10:
case 0x1c14:
case 0x1c19:
case 0x1c1f:
case 0x1c22:
case 0x1c23:
case 0x1c24:
case 0x1c27:
case 0x1c28:
case 0x1c29:
case 0x1c2b:
case 0x1c2d:
case 0x1c2f:
case 0x1c31:
case 0x1c34:
case 0x1c36:
case 0x1c39:
case 0x1c3b:
case 0x1c3e:
case 0x1c44:
case 0x1c46:
case 0x1c49:
case 0x1c4b:
case 0x1c4e:
case 0x1c51:
case 0x1c56:
case 0x1c58:
case 0x1c5b:
case 0x1c60:
case 0x1c62:
case 0x1c65:
case 0x1c69:
case 0x1c6d:
case 0x1c71:
case 0x1c74:
case 0x1c77:
case 0x1c7a:
case 0x1c7f:
case 0x1c81:
case 0x1c86:
case 0x1c88:
case 0x1c8b:
case 0x1cd0:
case 0x1cd5:
case 0x1cd7:
case 0x1cdc:
case 0x1cde:
case 0x1ce3:
case 0x1ce8:
case 0x1cec:
case 0x1ced:
case 0x1cf2:
case 0x1cf4:
case 0x1cf7:
case 0x1cfc:
case 0x1d01:
case 0x1d06:
case 0x1d0b:
case 0x1d11:
case 0x1d13:
case 0x1d16:
case 0x1d19:
case 0x1d1c:
case 0x1d1f:
case 0x1d20:
case 0x1d23:
case 0x1d26:
case 0x1d28:
case 0x1d2b:
case 0x1d30:
case 0x1d33:
case 0x1d37:
case 0x1d3b:
case 0x1d3e:
case 0x1d41:
case 0x1d43:
case 0x8b12:
case 0x8b13:
case 0x8b14:
case 0x8b15:
case 0x8b16:
case 0x8b17:
case 0x8b18:
case 0x8b1b:
case 0x8b1d:
case 0x8b20:
case 0x8b22:
case 0x8b25:
case 0x8b26:
case 0x8b2b:
case 0x8b30:
case 0x8b33:
case 0x8b37:
case 0x8b38:
case 0x8b39:
case 0x8b3a:
case 0x8b3b:
case 0x8b3c:
case 0x8b3d:
case 0x8b3e:
case 0x8b3f:
case 0x8b40:
case 0x8b41:
case 0x8b42:
case 0x8b43:
case 0x8b44:
case 0x8b45:
case 0x8b46:
case 0x8b49:
case 0x8b4b:
case 0x8b4e:
case 0x8b50:
case 0x8b52:
case 0x8b55:
case 0x8b56:
case 0x8b57:
case 0x8b58:
case 0x8b59:
case 0x8b5a:
case 0x8b5b:
case 0x8b5c:
case 0x8cae:
case 0x8caf:
case 0x8cb2:
case 0x8cb4:
case 0x8cb6:
case 0x8cb9:
case 0x8cbb:
case 0x8cbe:
case 0x8cc0:
case 0x8cc2:
case 0x8cc4:
case 0x8cc6:
case 0x8cc9:
case 0x8ccc:
case 0x8ccf:
case 0x8cd2:
case 0x8cd5:
case 0x8cd9:
case 0x8cdb:
case 0x8cde:
case 0x8ce1:
case 0x8ce4:
case 0x8ce7:
case 0x8ceb:
case 0x8ced:
case 0x8cf0:
case 0x8cf3:
case 0x8cf4:
case 0x8cf6:
case 0x8cf8:
case 0x8cfe:
case 0x8dc8:
case 0x8dc9:
case 0x8dca:
case 0x8dcb:
case 0x8dcc:
case 0x8dcd:
case 0x8dce:
case 0x8dd1:
case 0x8dd3:
case 0x8dd6:
case 0x8dd9:
case 0x8ddc:
case 0x8ddf:
case 0x8de2:
case 0x8de4:
case 0x8de5:
case 0x8de6:
case 0x8de8:
case 0x8dea:
case 0x8ded:
case 0x8df0:
case 0x8df1:
case 0x8df2:
case 0x8df5:
case 0x8df8:
case 0x8dfa:
case 0x8dfd:
case 0x8e00:
case 0x8e03:
case 0x8e06:
case 0x8e08:
case 0x8e0b:
case 0x8e0e:
case 0x8e11:
case 0x8e15:
case 0x8e19:
case 0x8e1c:
case 0x8e1e:
case 0x8e21:
case 0x8e23:
case 0x8e27:
case 0x8e29:
case 0x8e2c:
case 0x8e2e:
case 0x8e31:
case 0x8e33:
case 0x8e36:
case 0x8e38:
case 0x8e3a:
case 0x8e3c:
case 0x8e3f:
case 0x8e41:
case 0x8e43:
case 0x8e44:
case 0x8e47:
case 0x8e48:
case 0x8e49:
case 0x8e4a:
case 0x8e4b:
case 0x8e4e:
case 0x8e51:
case 0x8e54:
case 0x8e55:
case 0x8e56:
case 0x8e59:
case 0x8e5a:
case 0x8e5b:
case 0x8e5c:
case 0x8e5d:
case 0x8e62:
case 0x8e64:
case 0x8e67:
case 0x8e6a:
case 0x8e6d:
case 0x8e70:
case 0x8e73:
case 0x8e75:
case 0x8e78:
case 0x8e7a:
case 0x8e7c:
case 0x8e7e:
case 0x8e82:
case 0x8e84:
case 0x8e8a:
case 0x8e8c:
case 0x8e8e:
case 0x8e8f:
case 0x8e91:
case 0x8e94:
case 0x8e95:
case 0x8e96:
case 0x8e97:
case 0x8e9a:
case 0x8e9b:
case 0x8e9c:
case 0x8e9d:
case 0x8e9e:
case 0x8e9f:
case 0x8ea0:
case 0x8ea1:
case 0x8ea2:
case 0x8ea4:
case 0x8ea7:
case 0x8ea9:
case 0x8eab:
case 0x8ead:
case 0x8eaf:
case 0x8eb1:
case 0x8eb5:
case 0x8eb8:
case 0x8ebb:
case 0x8ebf:
case 0x8ec2:
case 0x8ec5:
case 0x8ec8:
case 0x8ecb:
case 0x8ece:
case 0x8ed0:
case 0x8ed2:
case 0x8ed4:
case 0x8ed6:
case 0x8ed8:
case 0x8edc:
case 0x8ede:
case 0x8ee0:
case 0x8ee3:
case 0x8ee6:
case 0x8eea:
case 0x8eed:
case 0x8ef0:
case 0x8ef3:
case 0x8ef6:
case 0x8ef9:
case 0x8efb:
case 0x8efd:
case 0x8eff:
case 0x8f01:
case 0x8f04:
case 0x8f06:
case 0x8f09:
case 0x8f0c:
case 0x8f0f:
case 0x8f12:
case 0x8f14:
case 0x8f17:
case 0x8f19:
case 0x8f1b:
case 0x8f1e:
case 0x8f21:
case 0x8f24:
case 0x8f27:
case 0x8f29:
case 0x8f2b:
case 0x8f2e:
case 0x8f31:
case 0x8f34:
case 0x8f35:
case 0x8f36:
case 0x8f37:
case 0x8f39:
case 0x8f3c:
case 0x8f3f:
case 0x8f42:
case 0x8f45:
case 0x8f47:
case 0x8f4a:
case 0x8f4d:
case 0x8f50:
case 0x8f52:
case 0x8f55:
case 0x8f58:
case 0x8f5b:
case 0x8f5e:
case 0x8f60:
case 0x8f63:
case 0x8f66:
case 0x8f69:
case 0x8f6c:
case 0x8f6d:
case 0x8f6f:
case 0x8f72:
case 0x8f75:
case 0x8f78:
case 0x8f7b:
case 0x8f7c:
case 0x8f7f:
case 0x8f82:
case 0x8f87:
case 0x8f8a:
case 0x8f8d:
case 0x8f90:
case 0x8f92:
case 0x8f93:
case 0x8f94:
case 0x8f96:
case 0x8fc9:
case 0x8fca:
case 0x8fcb:
case 0x8fcc:
case 0x8fce:
case 0x8fd1:
case 0x8fd3:
case 0x8fd7:
case 0x8fda:
case 0x8fdd:
case 0x8fdf:
case 0x8fe1:
case 0x8fe4:
case 0x8fe9:
case 0x8feb:
case 0x8fed:
case 0x8fef:
case 0x8ff1:
case 0x8ff3:
case 0x8ff8:
case 0x8ffd:
case 0x8fff:
case 0x9004:
case 0x9007:
case 0x900a:
case 0x900d:
case 0x900f:
case 0x9012:
case 0x9015:
case 0x9018:
case 0x901b:
case 0x901e:
case 0x9020:
case 0x9021:
case 0x9022:
case 0x9025:
case 0x9026:
case 0x9027:
case 0x9029:
case 0x902b:
case 0x902d:
case 0x902f:
case 0x9031:
case 0x9032:
case 0x9034:
case 0x9037:
case 0x903c:
case 0x903f:
case 0x9041:
case 0x9044:
case 0x9046:
case 0x9047:
case 0x904b:
case 0x904c:
case 0x904f:
case 0x9050:
case 0x9055:
case 0x9057:
case 0x9059:
case 0x905a:
case 0x905b:
case 0x905c:
case 0x905d:
case 0x905e:
case 0x9064:
case 0x9066:
case 0x9068:
case 0x9069:
case 0x906a:
case 0x906d:
case 0x906e:
case 0x9071:
case 0x9073:
case 0x9074:
case 0x9075:
case 0x907b:
case 0x907e:
case 0x907f:
case 0x9081:
case 0x9084:
case 0x9087:
case 0x9089:
case 0x908c:
case 0x908d:
case 0x908f:
case 0x9092:
case 0x9095:
case 0x9098:
case 0x909a:
case 0x909d:
case 0x909f:
case 0x90a2:
case 0x90a3:
case 0x90a7:
case 0x90a9:
case 0x90ab:
case 0x90ad:
case 0x90af:
case 0x90b1:
case 0x90b3:
case 0x90b5:
case 0x90b7:
case 0x90bc:
case 0x90bf:
case 0x90c0:
case 0x90c4:
case 0x90c6:
case 0x90ca:
case 0x90cc:
case 0x90d0:
case 0x90d2:
case 0x90d5:
case 0x90d7:
case 0x90d8:
case 0x90da:
case 0x90dc:
case 0x90de:
case 0x90e1:
case 0x90e3:
case 0x90e5:
case 0x90e7:
case 0x90e9:
case 0x90eb:
case 0x90ee:
case 0x90f0:
case 0x90f2:
case 0x90f4:
case 0x90f6:
case 0x90f8:
case 0x90fa:
case 0x90fc:
case 0x90fe:
case 0x9100:
case 0x9102:
case 0x9104:
case 0x9106:
case 0x9109:
case 0x910a:
case 0x910b:
case 0x9110:
case 0x9112:
case 0x9115:
case 0x9117:
case 0x911a:
case 0x911b:
case 0x911c:
case 0x911e:
case 0x9122:
case 0x9125:
case 0x9127:
case 0x9128:
case 0x9129:
case 0x912c:
case 0x912d:
case 0x9130:
case 0x9132:
case 0x9134:
case 0x9136:
case 0x9138:
case 0x913b:
case 0x913d:
case 0x9140:
case 0x9143:
case 0x9144:
case 0x9147:
case 0x9149:
case 0x914b:
case 0x914d:
case 0x914f:
case 0x9152:
case 0x9154:
case 0x9157:
case 0x915a:
case 0x915b:
case 0x915c:
case 0x915f:
case 0x9162:
case 0x9165:
case 0x9168:
case 0x916b:
case 0x916f:
case 0x9171:
case 0x9174:
case 0x9176:
case 0x9177:
case 0x9178:
case 0x917a:
case 0x917b:
case 0x917c:
case 0x917e:
case 0x9181:
case 0x9184:
case 0x9187:
case 0x918a:
case 0x918c:
case 0x918f:
case 0x9192:
case 0x9195:
case 0x9197:
case 0x919a:
case 0x919d:
case 0x91a0:
case 0x91a3:
case 0x91a5:
case 0x91a8:
case 0x91ab:
case 0x91ae:
case 0x91b1:
case 0x91b3:
case 0x91b6:
case 0x91b9:
case 0x91bc:
case 0x91bf:
case 0x91c1:
case 0x91c4:
case 0x91c7:
case 0x91ca:
case 0x91cd:
case 0x91cf:
case 0x91d2:
case 0x91d5:
case 0x91d8:
case 0x91db:
case 0x91dd:
case 0x91e0:
case 0x91e3:
case 0x91e6:
case 0x91e9:
case 0x91eb:
case 0x91ee:
case 0x91f1:
case 0x91f4:
case 0x91f7:
case 0x91f9:
case 0x91fc:
case 0x91ff:
case 0x9202:
case 0x9205:
case 0x9207:
case 0x920a:
case 0x920d:
case 0x9210:
case 0x9213:
case 0x9214:
case 0x9216:
case 0x9219:
case 0x921c:
case 0x921f:
case 0x9222:
case 0x9223:
case 0x9226:
case 0x922b:
case 0x922c:
case 0x922d:
case 0x9230:
case 0x9233:
case 0x9236:
case 0x9239:
case 0x923c:
case 0x923f:
case 0x9242:
case 0x9245:
case 0x9248:
case 0x924b:
case 0x924c:
case 0x924e:
case 0x9251:
case 0x9252:
case 0x9253:
case 0x9254:
case 0x9257:
case 0x9259:
case 0x925c:
case 0x925e:
case 0x9260:
case 0x9262:
case 0x9264:
case 0x9266:
case 0x9269:
case 0x926c:
case 0x926e:
case 0x9271:
case 0x9272:
case 0x9273:
case 0x9274:
case 0x9277:
case 0x9278:
case 0x9279:
case 0x927c:
case 0x927e:
case 0x927f:
case 0x9282:
case 0x9284:
case 0x9289:
case 0x928a:
case 0x928d:
case 0x9292:
case 0x9293:
case 0x9294:
case 0x9299:
case 0x929c:
case 0x929e:
case 0x92a0:
case 0x92a3:
case 0x92a5:
case 0x92a6:
case 0x92a7:
case 0x92a8:
case 0x92a9:
case 0x92ab:
case 0x92ae:
case 0x92b0:
case 0x92b3:
case 0x92b6:
case 0x92b8:
case 0x92bb:
case 0x92be:
case 0x92bf:
case 0x92c0:
case 0x92c1:
case 0x92c4:
case 0x92c7:
case 0x92c8:
case 0x92cb:
case 0x92cd:
case 0x92cf:
case 0x92d2:
case 0x92d4:
case 0x92d6:
case 0x92d7:
case 0x92d8:
case 0x92db:
case 0x92dd:
case 0x92e2:
case 0x92e3:
case 0x92e5:
case 0x92e8:
case 0x92ea:
case 0x92ed:
case 0x92ee:
case 0x92f0:
case 0x92f2:
case 0x92f4:
case 0x92f5:
case 0x92f8:
case 0x92fb:
case 0x92fe:
case 0x9300:
case 0x9303:
case 0x9305:
case 0x9306:
case 0x9307:
case 0x9308:
case 0xeb6c:
case 0xeb6f:
case 0xeb71:
case 0xeb74:
case 0xeb75:
case 0xeb78:
case 0xeb7a:
case 0xeb7b:
case 0xeb7d:
case 0xeb7e:
case 0xeb81:
case 0xeb83:
case 0xeb84:
case 0xeb85:
case 0xeb87:
case 0xeb89:
case 0xeb8b:
case 0xeb8c:
case 0xeb8e:
case 0xeb8f:
case 0xeb91:
case 0xeb92:
case 0xeb93:
case 0xeb94:
case 0xeb95:
case 0xeb96:
case 0xeb99:
case 0xeb9b:
case 0xeb9e:
case 0xeb9f:
case 0xeba2:
case 0xeba4:
case 0xeba5:
case 0xeba7:
case 0xeba8:
case 0xebaa:
case 0xebab:
case 0xebae:
case 0xebb0:
case 0xebb1:
case 0xebb2:
case 0xebb4:
case 0xebb5:
case 0xebb6:
case 0xebb8:
case 0xebb9:
case 0xebba:
case 0xebbc:
case 0xebbd:
case 0xebbe:
case 0xebc0:
case 0xebc1:
case 0xebc2:
case 0xebc4:
case 0xebc5:
case 0xebc6:
case 0xebc8:
case 0xebc9:
case 0xebca:
case 0xebcc:
case 0xebcd:
case 0xebcf:
case 0xebd2:
case 0xebd3:
case 0xebd6:
case 0xebd8:
case 0xebd9:
case 0xebda:
case 0xebdb:
return 1;default:return 0;}}
static void su_graph(KiMachine16 *m,uint16_t entry,const KiEconomyHooks *h,int su_skip_first_enter) {
const uint16_t su_core_cs=m->cs;int su_invoke_old=0;m->ip=entry;
SU_DISPATCH:
if(m->cs!=su_core_cs) {assert(!su_invoke_old);return;}
switch(m->ip) {
case 0x0:goto L_10000;
case 0x2:goto L_10002;
case 0x4:goto L_10004;
case 0x6:goto L_10006;
case 0x9:goto L_10009;
case 0xd:goto L_1000D;
case 0xf:goto L_1000F;
case 0x11:goto L_10011;
case 0x13:goto L_10013;
case 0x16:goto L_10016;
case 0x19:goto L_10019;
case 0x1e:goto L_1001E;
case 0x20:goto L_10020;
case 0x22:goto L_10022;
case 0x23:goto L_10023;
case 0x24:goto L_10024;
case 0x27:goto L_10027;
case 0x29:goto L_10029;
case 0x2b:goto L_1002B;
case 0x2c:goto L_1002C;
case 0x2d:goto L_1002D;
case 0x2f:goto L_1002F;
case 0x31:goto L_10031;
case 0x34:goto L_10034;
case 0x37:goto L_10037;
case 0x3a:goto L_1003A;
case 0x3d:goto L_1003D;
case 0x40:goto L_10040;
case 0x43:goto L_10043;
case 0x46:goto L_10046;
case 0x47:goto L_10047;
case 0x49:goto L_10049;
case 0x4b:goto L_1004B;
case 0x4e:goto L_1004E;
case 0x51:goto L_10051;
case 0x54:goto L_10054;
case 0x56:goto L_10056;
case 0x57:goto L_10057;
case 0x59:goto L_10059;
case 0x5b:goto L_1005B;
case 0x5e:goto L_1005E;
case 0x61:goto L_10061;
case 0x64:goto L_10064;
case 0x67:goto L_10067;
case 0x6a:goto L_1006A;
case 0x6b:goto L_1006B;
case 0x6c:goto L_1006C;
case 0x6f:goto L_1006F;
case 0x71:goto L_10071;
case 0x74:goto L_10074;
case 0x77:goto L_10077;
case 0x7a:goto L_1007A;
case 0x7d:goto L_1007D;
case 0x7f:goto L_1007F;
case 0x82:goto L_10082;
case 0x84:goto L_10084;
case 0x86:goto L_10086;
case 0x89:goto L_10089;
case 0x8c:goto L_1008C;
case 0x8f:goto L_1008F;
case 0x90:goto L_10090;
case 0x92:goto L_10092;
case 0x95:goto L_10095;
case 0x99:goto L_10099;
case 0x9b:goto L_1009B;
case 0x9e:goto L_1009E;
case 0xa1:goto L_100A1;
case 0xa5:goto L_100A5;
case 0xa7:goto L_100A7;
case 0xa9:goto L_100A9;
case 0xac:goto L_100AC;
case 0xaf:goto L_100AF;
case 0xb2:goto L_100B2;
case 0xb5:goto L_100B5;
case 0xb9:goto L_100B9;
case 0xbb:goto L_100BB;
case 0xbd:goto L_100BD;
case 0xc0:goto L_100C0;
case 0xc3:goto L_100C3;
case 0xc6:goto L_100C6;
case 0xc9:goto L_100C9;
case 0xcd:goto L_100CD;
case 0xcf:goto L_100CF;
case 0xd2:goto L_100D2;
case 0xd5:goto L_100D5;
case 0xd9:goto L_100D9;
case 0xdb:goto L_100DB;
case 0xde:goto L_100DE;
case 0x155:goto L_10155;
case 0x158:goto L_10158;
case 0x15d:goto L_1015D;
case 0x160:goto L_10160;
case 0x162:goto L_10162;
case 0x167:goto L_10167;
case 0x16a:goto L_1016A;
case 0x16d:goto L_1016D;
case 0x16f:goto L_1016F;
case 0x174:goto L_10174;
case 0x177:goto L_10177;
case 0x179:goto L_10179;
case 0x17e:goto L_1017E;
case 0x181:goto L_10181;
case 0x184:goto L_10184;
case 0x189:goto L_10189;
case 0x18c:goto L_1018C;
case 0x18f:goto L_1018F;
case 0x192:goto L_10192;
case 0x197:goto L_10197;
case 0x19a:goto L_1019A;
case 0x19d:goto L_1019D;
case 0x1a0:goto L_101A0;
case 0x1a5:goto L_101A5;
case 0x1a8:goto L_101A8;
case 0x1ab:goto L_101AB;
case 0x1ae:goto L_101AE;
case 0x1b3:goto L_101B3;
case 0x210:goto L_10210;
case 0x211:goto L_10211;
case 0x212:goto L_10212;
case 0x213:goto L_10213;
case 0x215:goto L_10215;
case 0x217:goto L_10217;
case 0x219:goto L_10219;
case 0x21e:goto L_1021E;
case 0x220:goto L_10220;
case 0x222:goto L_10222;
case 0x224:goto L_10224;
case 0x226:goto L_10226;
case 0x22c:goto L_1022C;
case 0x232:goto L_10232;
case 0x233:goto L_10233;
case 0x234:goto L_10234;
case 0x235:goto L_10235;
case 0x236:goto L_10236;
case 0x237:goto L_10237;
case 0x239:goto L_10239;
case 0x23b:goto L_1023B;
case 0x23d:goto L_1023D;
case 0x23f:goto L_1023F;
case 0x240:goto L_10240;
case 0x33b:goto L_1033B;
case 0x33c:goto L_1033C;
case 0x33d:goto L_1033D;
case 0x33e:goto L_1033E;
case 0x340:goto L_10340;
case 0x342:goto L_10342;
case 0x345:goto L_10345;
case 0x348:goto L_10348;
case 0x34a:goto L_1034A;
case 0x34b:goto L_1034B;
case 0x34c:goto L_1034C;
case 0x34d:goto L_1034D;
case 0x34e:goto L_1034E;
case 0x34f:goto L_1034F;
case 0x352:goto L_10352;
case 0x354:goto L_10354;
case 0x355:goto L_10355;
case 0x356:goto L_10356;
case 0x357:goto L_10357;
case 0x358:goto L_10358;
case 0x35e:goto L_1035E;
case 0x364:goto L_10364;
case 0x366:goto L_10366;
case 0x36b:goto L_1036B;
case 0x36c:goto L_1036C;
case 0x612:goto L_10612;
case 0x614:goto L_10614;
case 0x616:goto L_10616;
case 0x618:goto L_10618;
case 0x61a:goto L_1061A;
case 0x61c:goto L_1061C;
case 0x61e:goto L_1061E;
case 0x1a6e:goto L_11A6E;
case 0x1a70:goto L_11A70;
case 0x1a72:goto L_11A72;
case 0x1a74:goto L_11A74;
case 0x1a76:goto L_11A76;
case 0x1a79:goto L_11A79;
case 0x1a7b:goto L_11A7B;
case 0x1a7e:goto L_11A7E;
case 0x1a84:goto L_11A84;
case 0x1a87:goto L_11A87;
case 0x1a8a:goto L_11A8A;
case 0x1a8d:goto L_11A8D;
case 0x1a90:goto L_11A90;
case 0x1a93:goto L_11A93;
case 0x1a96:goto L_11A96;
case 0x1a99:goto L_11A99;
case 0x1a9c:goto L_11A9C;
case 0x1a9e:goto L_11A9E;
case 0x1aa1:goto L_11AA1;
case 0x1aa4:goto L_11AA4;
case 0x1aa7:goto L_11AA7;
case 0x1aad:goto L_11AAD;
case 0x1aae:goto L_11AAE;
case 0x1ab1:goto L_11AB1;
case 0x1ab4:goto L_11AB4;
case 0x1ab7:goto L_11AB7;
case 0x1aba:goto L_11ABA;
case 0x1abc:goto L_11ABC;
case 0x1abf:goto L_11ABF;
case 0x1ac2:goto L_11AC2;
case 0x1ac3:goto L_11AC3;
case 0x1ac6:goto L_11AC6;
case 0x1ac9:goto L_11AC9;
case 0x1acc:goto L_11ACC;
case 0x1acf:goto L_11ACF;
case 0x1ad2:goto L_11AD2;
case 0x1ad5:goto L_11AD5;
case 0x1ad7:goto L_11AD7;
case 0x1ad9:goto L_11AD9;
case 0x1adc:goto L_11ADC;
case 0x1ade:goto L_11ADE;
case 0x1ae0:goto L_11AE0;
case 0x1ae3:goto L_11AE3;
case 0x1ae5:goto L_11AE5;
case 0x1ae8:goto L_11AE8;
case 0x1aea:goto L_11AEA;
case 0x1aed:goto L_11AED;
case 0x1aef:goto L_11AEF;
case 0x1af1:goto L_11AF1;
case 0x1af3:goto L_11AF3;
case 0x1af6:goto L_11AF6;
case 0x1af8:goto L_11AF8;
case 0x1afd:goto L_11AFD;
case 0x1b00:goto L_11B00;
case 0x1b03:goto L_11B03;
case 0x1b05:goto L_11B05;
case 0x1b08:goto L_11B08;
case 0x1b0a:goto L_11B0A;
case 0x1b0c:goto L_11B0C;
case 0x1b0e:goto L_11B0E;
case 0x1b10:goto L_11B10;
case 0x1b12:goto L_11B12;
case 0x1b17:goto L_11B17;
case 0x1b1c:goto L_11B1C;
case 0x1b1e:goto L_11B1E;
case 0x1b20:goto L_11B20;
case 0x1b25:goto L_11B25;
case 0x1b27:goto L_11B27;
case 0x1b29:goto L_11B29;
case 0x1b2c:goto L_11B2C;
case 0x1b2f:goto L_11B2F;
case 0x1b34:goto L_11B34;
case 0x1b39:goto L_11B39;
case 0x1b3c:goto L_11B3C;
case 0x1b3e:goto L_11B3E;
case 0x1b40:goto L_11B40;
case 0x1b42:goto L_11B42;
case 0x1b44:goto L_11B44;
case 0x1b48:goto L_11B48;
case 0x1b4c:goto L_11B4C;
case 0x1b4f:goto L_11B4F;
case 0x1b52:goto L_11B52;
case 0x1b55:goto L_11B55;
case 0x1b57:goto L_11B57;
case 0x1b59:goto L_11B59;
case 0x1be0:goto L_11BE0;
case 0x1be2:goto L_11BE2;
case 0x1be4:goto L_11BE4;
case 0x1be6:goto L_11BE6;
case 0x1be9:goto L_11BE9;
case 0x1bec:goto L_11BEC;
case 0x1bf1:goto L_11BF1;
case 0x1bf6:goto L_11BF6;
case 0x1bf9:goto L_11BF9;
case 0x1bfd:goto L_11BFD;
case 0x1c01:goto L_11C01;
case 0x1c04:goto L_11C04;
case 0x1c07:goto L_11C07;
case 0x1c0a:goto L_11C0A;
case 0x1c0d:goto L_11C0D;
case 0x1c10:goto L_11C10;
case 0x1c14:goto L_11C14;
case 0x1c19:goto L_11C19;
case 0x1c1f:goto L_11C1F;
case 0x1c22:goto L_11C22;
case 0x1c23:goto L_11C23;
case 0x1c24:goto L_11C24;
case 0x1c27:goto L_11C27;
case 0x1c28:goto L_11C28;
case 0x1c29:goto L_11C29;
case 0x1c2b:goto L_11C2B;
case 0x1c2d:goto L_11C2D;
case 0x1c2f:goto L_11C2F;
case 0x1c31:goto L_11C31;
case 0x1c34:goto L_11C34;
case 0x1c36:goto L_11C36;
case 0x1c39:goto L_11C39;
case 0x1c3b:goto L_11C3B;
case 0x1c3e:goto L_11C3E;
case 0x1c44:goto L_11C44;
case 0x1c46:goto L_11C46;
case 0x1c49:goto L_11C49;
case 0x1c4b:goto L_11C4B;
case 0x1c4e:goto L_11C4E;
case 0x1c51:goto L_11C51;
case 0x1c56:goto L_11C56;
case 0x1c58:goto L_11C58;
case 0x1c5b:goto L_11C5B;
case 0x1c60:goto L_11C60;
case 0x1c62:goto L_11C62;
case 0x1c65:goto L_11C65;
case 0x1c69:goto L_11C69;
case 0x1c6d:goto L_11C6D;
case 0x1c71:goto L_11C71;
case 0x1c74:goto L_11C74;
case 0x1c77:goto L_11C77;
case 0x1c7a:goto L_11C7A;
case 0x1c7f:goto L_11C7F;
case 0x1c81:goto L_11C81;
case 0x1c86:goto L_11C86;
case 0x1c88:goto L_11C88;
case 0x1c8b:goto L_11C8B;
case 0x1cd0:goto L_11CD0;
case 0x1cd5:goto L_11CD5;
case 0x1cd7:goto L_11CD7;
case 0x1cdc:goto L_11CDC;
case 0x1cde:goto L_11CDE;
case 0x1ce3:goto L_11CE3;
case 0x1ce8:goto L_11CE8;
case 0x1cec:goto L_11CEC;
case 0x1ced:goto L_11CED;
case 0x1cf2:goto L_11CF2;
case 0x1cf4:goto L_11CF4;
case 0x1cf7:goto L_11CF7;
case 0x1cfc:goto L_11CFC;
case 0x1d01:goto L_11D01;
case 0x1d06:goto L_11D06;
case 0x1d0b:goto L_11D0B;
case 0x1d11:goto L_11D11;
case 0x1d13:goto L_11D13;
case 0x1d16:goto L_11D16;
case 0x1d19:goto L_11D19;
case 0x1d1c:goto L_11D1C;
case 0x1d1f:goto L_11D1F;
case 0x1d20:goto L_11D20;
case 0x1d23:goto L_11D23;
case 0x1d26:goto L_11D26;
case 0x1d28:goto L_11D28;
case 0x1d2b:goto L_11D2B;
case 0x1d30:goto L_11D30;
case 0x1d33:goto L_11D33;
case 0x1d37:goto L_11D37;
case 0x1d3b:goto L_11D3B;
case 0x1d3e:goto L_11D3E;
case 0x1d41:goto L_11D41;
case 0x1d43:goto L_11D43;
case 0x8b12:goto L_18B12;
case 0x8b13:goto L_18B13;
case 0x8b14:goto L_18B14;
case 0x8b15:goto L_18B15;
case 0x8b16:goto L_18B16;
case 0x8b17:goto L_18B17;
case 0x8b18:goto L_18B18;
case 0x8b1b:goto L_18B1B;
case 0x8b1d:goto L_18B1D;
case 0x8b20:goto L_18B20;
case 0x8b22:goto L_18B22;
case 0x8b25:goto L_18B25;
case 0x8b26:goto L_18B26;
case 0x8b2b:goto L_18B2B;
case 0x8b30:goto L_18B30;
case 0x8b33:goto L_18B33;
case 0x8b37:goto L_18B37;
case 0x8b38:goto L_18B38;
case 0x8b39:goto L_18B39;
case 0x8b3a:goto L_18B3A;
case 0x8b3b:goto L_18B3B;
case 0x8b3c:goto L_18B3C;
case 0x8b3d:goto L_18B3D;
case 0x8b3e:goto L_18B3E;
case 0x8b3f:goto L_18B3F;
case 0x8b40:goto L_18B40;
case 0x8b41:goto L_18B41;
case 0x8b42:goto L_18B42;
case 0x8b43:goto L_18B43;
case 0x8b44:goto L_18B44;
case 0x8b45:goto L_18B45;
case 0x8b46:goto L_18B46;
case 0x8b49:goto L_18B49;
case 0x8b4b:goto L_18B4B;
case 0x8b4e:goto L_18B4E;
case 0x8b50:goto L_18B50;
case 0x8b52:goto L_18B52;
case 0x8b55:goto L_18B55;
case 0x8b56:goto L_18B56;
case 0x8b57:goto L_18B57;
case 0x8b58:goto L_18B58;
case 0x8b59:goto L_18B59;
case 0x8b5a:goto L_18B5A;
case 0x8b5b:goto L_18B5B;
case 0x8b5c:goto L_18B5C;
case 0x8cae:goto L_18CAE;
case 0x8caf:goto L_18CAF;
case 0x8cb2:goto L_18CB2;
case 0x8cb4:goto L_18CB4;
case 0x8cb6:goto L_18CB6;
case 0x8cb9:goto L_18CB9;
case 0x8cbb:goto L_18CBB;
case 0x8cbe:goto L_18CBE;
case 0x8cc0:goto L_18CC0;
case 0x8cc2:goto L_18CC2;
case 0x8cc4:goto L_18CC4;
case 0x8cc6:goto L_18CC6;
case 0x8cc9:goto L_18CC9;
case 0x8ccc:goto L_18CCC;
case 0x8ccf:goto L_18CCF;
case 0x8cd2:goto L_18CD2;
case 0x8cd5:goto L_18CD5;
case 0x8cd9:goto L_18CD9;
case 0x8cdb:goto L_18CDB;
case 0x8cde:goto L_18CDE;
case 0x8ce1:goto L_18CE1;
case 0x8ce4:goto L_18CE4;
case 0x8ce7:goto L_18CE7;
case 0x8ceb:goto L_18CEB;
case 0x8ced:goto L_18CED;
case 0x8cf0:goto L_18CF0;
case 0x8cf3:goto L_18CF3;
case 0x8cf4:goto L_18CF4;
case 0x8cf6:goto L_18CF6;
case 0x8cf8:goto L_18CF8;
case 0x8cfe:goto L_18CFE;
case 0x8dc8:goto L_18DC8;
case 0x8dc9:goto L_18DC9;
case 0x8dca:goto L_18DCA;
case 0x8dcb:goto L_18DCB;
case 0x8dcc:goto L_18DCC;
case 0x8dcd:goto L_18DCD;
case 0x8dce:goto L_18DCE;
case 0x8dd1:goto L_18DD1;
case 0x8dd3:goto L_18DD3;
case 0x8dd6:goto L_18DD6;
case 0x8dd9:goto L_18DD9;
case 0x8ddc:goto L_18DDC;
case 0x8ddf:goto L_18DDF;
case 0x8de2:goto L_18DE2;
case 0x8de4:goto L_18DE4;
case 0x8de5:goto L_18DE5;
case 0x8de6:goto L_18DE6;
case 0x8de8:goto L_18DE8;
case 0x8dea:goto L_18DEA;
case 0x8ded:goto L_18DED;
case 0x8df0:goto L_18DF0;
case 0x8df1:goto L_18DF1;
case 0x8df2:goto L_18DF2;
case 0x8df5:goto L_18DF5;
case 0x8df8:goto L_18DF8;
case 0x8dfa:goto L_18DFA;
case 0x8dfd:goto L_18DFD;
case 0x8e00:goto L_18E00;
case 0x8e03:goto L_18E03;
case 0x8e06:goto L_18E06;
case 0x8e08:goto L_18E08;
case 0x8e0b:goto L_18E0B;
case 0x8e0e:goto L_18E0E;
case 0x8e11:goto L_18E11;
case 0x8e15:goto L_18E15;
case 0x8e19:goto L_18E19;
case 0x8e1c:goto L_18E1C;
case 0x8e1e:goto L_18E1E;
case 0x8e21:goto L_18E21;
case 0x8e23:goto L_18E23;
case 0x8e27:goto L_18E27;
case 0x8e29:goto L_18E29;
case 0x8e2c:goto L_18E2C;
case 0x8e2e:goto L_18E2E;
case 0x8e31:goto L_18E31;
case 0x8e33:goto L_18E33;
case 0x8e36:goto L_18E36;
case 0x8e38:goto L_18E38;
case 0x8e3a:goto L_18E3A;
case 0x8e3c:goto L_18E3C;
case 0x8e3f:goto L_18E3F;
case 0x8e41:goto L_18E41;
case 0x8e43:goto L_18E43;
case 0x8e44:goto L_18E44;
case 0x8e47:goto L_18E47;
case 0x8e48:goto L_18E48;
case 0x8e49:goto L_18E49;
case 0x8e4a:goto L_18E4A;
case 0x8e4b:goto L_18E4B;
case 0x8e4e:goto L_18E4E;
case 0x8e51:goto L_18E51;
case 0x8e54:goto L_18E54;
case 0x8e55:goto L_18E55;
case 0x8e56:goto L_18E56;
case 0x8e59:goto L_18E59;
case 0x8e5a:goto L_18E5A;
case 0x8e5b:goto L_18E5B;
case 0x8e5c:goto L_18E5C;
case 0x8e5d:goto L_18E5D;
case 0x8e62:goto L_18E62;
case 0x8e64:goto L_18E64;
case 0x8e67:goto L_18E67;
case 0x8e6a:goto L_18E6A;
case 0x8e6d:goto L_18E6D;
case 0x8e70:goto L_18E70;
case 0x8e73:goto L_18E73;
case 0x8e75:goto L_18E75;
case 0x8e78:goto L_18E78;
case 0x8e7a:goto L_18E7A;
case 0x8e7c:goto L_18E7C;
case 0x8e7e:goto L_18E7E;
case 0x8e82:goto L_18E82;
case 0x8e84:goto L_18E84;
case 0x8e8a:goto L_18E8A;
case 0x8e8c:goto L_18E8C;
case 0x8e8e:goto L_18E8E;
case 0x8e8f:goto L_18E8F;
case 0x8e91:goto L_18E91;
case 0x8e94:goto L_18E94;
case 0x8e95:goto L_18E95;
case 0x8e96:goto L_18E96;
case 0x8e97:goto L_18E97;
case 0x8e9a:goto L_18E9A;
case 0x8e9b:goto L_18E9B;
case 0x8e9c:goto L_18E9C;
case 0x8e9d:goto L_18E9D;
case 0x8e9e:goto L_18E9E;
case 0x8e9f:goto L_18E9F;
case 0x8ea0:goto L_18EA0;
case 0x8ea1:goto L_18EA1;
case 0x8ea2:goto L_18EA2;
case 0x8ea4:goto L_18EA4;
case 0x8ea7:goto L_18EA7;
case 0x8ea9:goto L_18EA9;
case 0x8eab:goto L_18EAB;
case 0x8ead:goto L_18EAD;
case 0x8eaf:goto L_18EAF;
case 0x8eb1:goto L_18EB1;
case 0x8eb5:goto L_18EB5;
case 0x8eb8:goto L_18EB8;
case 0x8ebb:goto L_18EBB;
case 0x8ebf:goto L_18EBF;
case 0x8ec2:goto L_18EC2;
case 0x8ec5:goto L_18EC5;
case 0x8ec8:goto L_18EC8;
case 0x8ecb:goto L_18ECB;
case 0x8ece:goto L_18ECE;
case 0x8ed0:goto L_18ED0;
case 0x8ed2:goto L_18ED2;
case 0x8ed4:goto L_18ED4;
case 0x8ed6:goto L_18ED6;
case 0x8ed8:goto L_18ED8;
case 0x8edc:goto L_18EDC;
case 0x8ede:goto L_18EDE;
case 0x8ee0:goto L_18EE0;
case 0x8ee3:goto L_18EE3;
case 0x8ee6:goto L_18EE6;
case 0x8eea:goto L_18EEA;
case 0x8eed:goto L_18EED;
case 0x8ef0:goto L_18EF0;
case 0x8ef3:goto L_18EF3;
case 0x8ef6:goto L_18EF6;
case 0x8ef9:goto L_18EF9;
case 0x8efb:goto L_18EFB;
case 0x8efd:goto L_18EFD;
case 0x8eff:goto L_18EFF;
case 0x8f01:goto L_18F01;
case 0x8f04:goto L_18F04;
case 0x8f06:goto L_18F06;
case 0x8f09:goto L_18F09;
case 0x8f0c:goto L_18F0C;
case 0x8f0f:goto L_18F0F;
case 0x8f12:goto L_18F12;
case 0x8f14:goto L_18F14;
case 0x8f17:goto L_18F17;
case 0x8f19:goto L_18F19;
case 0x8f1b:goto L_18F1B;
case 0x8f1e:goto L_18F1E;
case 0x8f21:goto L_18F21;
case 0x8f24:goto L_18F24;
case 0x8f27:goto L_18F27;
case 0x8f29:goto L_18F29;
case 0x8f2b:goto L_18F2B;
case 0x8f2e:goto L_18F2E;
case 0x8f31:goto L_18F31;
case 0x8f34:goto L_18F34;
case 0x8f35:goto L_18F35;
case 0x8f36:goto L_18F36;
case 0x8f37:goto L_18F37;
case 0x8f39:goto L_18F39;
case 0x8f3c:goto L_18F3C;
case 0x8f3f:goto L_18F3F;
case 0x8f42:goto L_18F42;
case 0x8f45:goto L_18F45;
case 0x8f47:goto L_18F47;
case 0x8f4a:goto L_18F4A;
case 0x8f4d:goto L_18F4D;
case 0x8f50:goto L_18F50;
case 0x8f52:goto L_18F52;
case 0x8f55:goto L_18F55;
case 0x8f58:goto L_18F58;
case 0x8f5b:goto L_18F5B;
case 0x8f5e:goto L_18F5E;
case 0x8f60:goto L_18F60;
case 0x8f63:goto L_18F63;
case 0x8f66:goto L_18F66;
case 0x8f69:goto L_18F69;
case 0x8f6c:goto L_18F6C;
case 0x8f6d:goto L_18F6D;
case 0x8f6f:goto L_18F6F;
case 0x8f72:goto L_18F72;
case 0x8f75:goto L_18F75;
case 0x8f78:goto L_18F78;
case 0x8f7b:goto L_18F7B;
case 0x8f7c:goto L_18F7C;
case 0x8f7f:goto L_18F7F;
case 0x8f82:goto L_18F82;
case 0x8f87:goto L_18F87;
case 0x8f8a:goto L_18F8A;
case 0x8f8d:goto L_18F8D;
case 0x8f90:goto L_18F90;
case 0x8f92:goto L_18F92;
case 0x8f93:goto L_18F93;
case 0x8f94:goto L_18F94;
case 0x8f96:goto L_18F96;
case 0x8fc9:goto L_18FC9;
case 0x8fca:goto L_18FCA;
case 0x8fcb:goto L_18FCB;
case 0x8fcc:goto L_18FCC;
case 0x8fce:goto L_18FCE;
case 0x8fd1:goto L_18FD1;
case 0x8fd3:goto L_18FD3;
case 0x8fd7:goto L_18FD7;
case 0x8fda:goto L_18FDA;
case 0x8fdd:goto L_18FDD;
case 0x8fdf:goto L_18FDF;
case 0x8fe1:goto L_18FE1;
case 0x8fe4:goto L_18FE4;
case 0x8fe9:goto L_18FE9;
case 0x8feb:goto L_18FEB;
case 0x8fed:goto L_18FED;
case 0x8fef:goto L_18FEF;
case 0x8ff1:goto L_18FF1;
case 0x8ff3:goto L_18FF3;
case 0x8ff8:goto L_18FF8;
case 0x8ffd:goto L_18FFD;
case 0x8fff:goto L_18FFF;
case 0x9004:goto L_19004;
case 0x9007:goto L_19007;
case 0x900a:goto L_1900A;
case 0x900d:goto L_1900D;
case 0x900f:goto L_1900F;
case 0x9012:goto L_19012;
case 0x9015:goto L_19015;
case 0x9018:goto L_19018;
case 0x901b:goto L_1901B;
case 0x901e:goto L_1901E;
case 0x9020:goto L_19020;
case 0x9021:goto L_19021;
case 0x9022:goto L_19022;
case 0x9025:goto L_19025;
case 0x9026:goto L_19026;
case 0x9027:goto L_19027;
case 0x9029:goto L_19029;
case 0x902b:goto L_1902B;
case 0x902d:goto L_1902D;
case 0x902f:goto L_1902F;
case 0x9031:goto L_19031;
case 0x9032:goto L_19032;
case 0x9034:goto L_19034;
case 0x9037:goto L_19037;
case 0x903c:goto L_1903C;
case 0x903f:goto L_1903F;
case 0x9041:goto L_19041;
case 0x9044:goto L_19044;
case 0x9046:goto L_19046;
case 0x9047:goto L_19047;
case 0x904b:goto L_1904B;
case 0x904c:goto L_1904C;
case 0x904f:goto L_1904F;
case 0x9050:goto L_19050;
case 0x9055:goto L_19055;
case 0x9057:goto L_19057;
case 0x9059:goto L_19059;
case 0x905a:goto L_1905A;
case 0x905b:goto L_1905B;
case 0x905c:goto L_1905C;
case 0x905d:goto L_1905D;
case 0x905e:goto L_1905E;
case 0x9064:goto L_19064;
case 0x9066:goto L_19066;
case 0x9068:goto L_19068;
case 0x9069:goto L_19069;
case 0x906a:goto L_1906A;
case 0x906d:goto L_1906D;
case 0x906e:goto L_1906E;
case 0x9071:goto L_19071;
case 0x9073:goto L_19073;
case 0x9074:goto L_19074;
case 0x9075:goto L_19075;
case 0x907b:goto L_1907B;
case 0x907e:goto L_1907E;
case 0x907f:goto L_1907F;
case 0x9081:goto L_19081;
case 0x9084:goto L_19084;
case 0x9087:goto L_19087;
case 0x9089:goto L_19089;
case 0x908c:goto L_1908C;
case 0x908d:goto L_1908D;
case 0x908f:goto L_1908F;
case 0x9092:goto L_19092;
case 0x9095:goto L_19095;
case 0x9098:goto L_19098;
case 0x909a:goto L_1909A;
case 0x909d:goto L_1909D;
case 0x909f:goto L_1909F;
case 0x90a2:goto L_190A2;
case 0x90a3:goto L_190A3;
case 0x90a7:goto L_190A7;
case 0x90a9:goto L_190A9;
case 0x90ab:goto L_190AB;
case 0x90ad:goto L_190AD;
case 0x90af:goto L_190AF;
case 0x90b1:goto L_190B1;
case 0x90b3:goto L_190B3;
case 0x90b5:goto L_190B5;
case 0x90b7:goto L_190B7;
case 0x90bc:goto L_190BC;
case 0x90bf:goto L_190BF;
case 0x90c0:goto L_190C0;
case 0x90c4:goto L_190C4;
case 0x90c6:goto L_190C6;
case 0x90ca:goto L_190CA;
case 0x90cc:goto L_190CC;
case 0x90d0:goto L_190D0;
case 0x90d2:goto L_190D2;
case 0x90d5:goto L_190D5;
case 0x90d7:goto L_190D7;
case 0x90d8:goto L_190D8;
case 0x90da:goto L_190DA;
case 0x90dc:goto L_190DC;
case 0x90de:goto L_190DE;
case 0x90e1:goto L_190E1;
case 0x90e3:goto L_190E3;
case 0x90e5:goto L_190E5;
case 0x90e7:goto L_190E7;
case 0x90e9:goto L_190E9;
case 0x90eb:goto L_190EB;
case 0x90ee:goto L_190EE;
case 0x90f0:goto L_190F0;
case 0x90f2:goto L_190F2;
case 0x90f4:goto L_190F4;
case 0x90f6:goto L_190F6;
case 0x90f8:goto L_190F8;
case 0x90fa:goto L_190FA;
case 0x90fc:goto L_190FC;
case 0x90fe:goto L_190FE;
case 0x9100:goto L_19100;
case 0x9102:goto L_19102;
case 0x9104:goto L_19104;
case 0x9106:goto L_19106;
case 0x9109:goto L_19109;
case 0x910a:goto L_1910A;
case 0x910b:goto L_1910B;
case 0x9110:goto L_19110;
case 0x9112:goto L_19112;
case 0x9115:goto L_19115;
case 0x9117:goto L_19117;
case 0x911a:goto L_1911A;
case 0x911b:goto L_1911B;
case 0x911c:goto L_1911C;
case 0x911e:goto L_1911E;
case 0x9122:goto L_19122;
case 0x9125:goto L_19125;
case 0x9127:goto L_19127;
case 0x9128:goto L_19128;
case 0x9129:goto L_19129;
case 0x912c:goto L_1912C;
case 0x912d:goto L_1912D;
case 0x9130:goto L_19130;
case 0x9132:goto L_19132;
case 0x9134:goto L_19134;
case 0x9136:goto L_19136;
case 0x9138:goto L_19138;
case 0x913b:goto L_1913B;
case 0x913d:goto L_1913D;
case 0x9140:goto L_19140;
case 0x9143:goto L_19143;
case 0x9144:goto L_19144;
case 0x9147:goto L_19147;
case 0x9149:goto L_19149;
case 0x914b:goto L_1914B;
case 0x914d:goto L_1914D;
case 0x914f:goto L_1914F;
case 0x9152:goto L_19152;
case 0x9154:goto L_19154;
case 0x9157:goto L_19157;
case 0x915a:goto L_1915A;
case 0x915b:goto L_1915B;
case 0x915c:goto L_1915C;
case 0x915f:goto L_1915F;
case 0x9162:goto L_19162;
case 0x9165:goto L_19165;
case 0x9168:goto L_19168;
case 0x916b:goto L_1916B;
case 0x916f:goto L_1916F;
case 0x9171:goto L_19171;
case 0x9174:goto L_19174;
case 0x9176:goto L_19176;
case 0x9177:goto L_19177;
case 0x9178:goto L_19178;
case 0x917a:goto L_1917A;
case 0x917b:goto L_1917B;
case 0x917c:goto L_1917C;
case 0x917e:goto L_1917E;
case 0x9181:goto L_19181;
case 0x9184:goto L_19184;
case 0x9187:goto L_19187;
case 0x918a:goto L_1918A;
case 0x918c:goto L_1918C;
case 0x918f:goto L_1918F;
case 0x9192:goto L_19192;
case 0x9195:goto L_19195;
case 0x9197:goto L_19197;
case 0x919a:goto L_1919A;
case 0x919d:goto L_1919D;
case 0x91a0:goto L_191A0;
case 0x91a3:goto L_191A3;
case 0x91a5:goto L_191A5;
case 0x91a8:goto L_191A8;
case 0x91ab:goto L_191AB;
case 0x91ae:goto L_191AE;
case 0x91b1:goto L_191B1;
case 0x91b3:goto L_191B3;
case 0x91b6:goto L_191B6;
case 0x91b9:goto L_191B9;
case 0x91bc:goto L_191BC;
case 0x91bf:goto L_191BF;
case 0x91c1:goto L_191C1;
case 0x91c4:goto L_191C4;
case 0x91c7:goto L_191C7;
case 0x91ca:goto L_191CA;
case 0x91cd:goto L_191CD;
case 0x91cf:goto L_191CF;
case 0x91d2:goto L_191D2;
case 0x91d5:goto L_191D5;
case 0x91d8:goto L_191D8;
case 0x91db:goto L_191DB;
case 0x91dd:goto L_191DD;
case 0x91e0:goto L_191E0;
case 0x91e3:goto L_191E3;
case 0x91e6:goto L_191E6;
case 0x91e9:goto L_191E9;
case 0x91eb:goto L_191EB;
case 0x91ee:goto L_191EE;
case 0x91f1:goto L_191F1;
case 0x91f4:goto L_191F4;
case 0x91f7:goto L_191F7;
case 0x91f9:goto L_191F9;
case 0x91fc:goto L_191FC;
case 0x91ff:goto L_191FF;
case 0x9202:goto L_19202;
case 0x9205:goto L_19205;
case 0x9207:goto L_19207;
case 0x920a:goto L_1920A;
case 0x920d:goto L_1920D;
case 0x9210:goto L_19210;
case 0x9213:goto L_19213;
case 0x9214:goto L_19214;
case 0x9216:goto L_19216;
case 0x9219:goto L_19219;
case 0x921c:goto L_1921C;
case 0x921f:goto L_1921F;
case 0x9222:goto L_19222;
case 0x9223:goto L_19223;
case 0x9226:goto L_19226;
case 0x922b:goto L_1922B;
case 0x922c:goto L_1922C;
case 0x922d:goto L_1922D;
case 0x9230:goto L_19230;
case 0x9233:goto L_19233;
case 0x9236:goto L_19236;
case 0x9239:goto L_19239;
case 0x923c:goto L_1923C;
case 0x923f:goto L_1923F;
case 0x9242:goto L_19242;
case 0x9245:goto L_19245;
case 0x9248:goto L_19248;
case 0x924b:goto L_1924B;
case 0x924c:goto L_1924C;
case 0x924e:goto L_1924E;
case 0x9251:goto L_19251;
case 0x9252:goto L_19252;
case 0x9253:goto L_19253;
case 0x9254:goto L_19254;
case 0x9257:goto L_19257;
case 0x9259:goto L_19259;
case 0x925c:goto L_1925C;
case 0x925e:goto L_1925E;
case 0x9260:goto L_19260;
case 0x9262:goto L_19262;
case 0x9264:goto L_19264;
case 0x9266:goto L_19266;
case 0x9269:goto L_19269;
case 0x926c:goto L_1926C;
case 0x926e:goto L_1926E;
case 0x9271:goto L_19271;
case 0x9272:goto L_19272;
case 0x9273:goto L_19273;
case 0x9274:goto L_19274;
case 0x9277:goto L_19277;
case 0x9278:goto L_19278;
case 0x9279:goto L_19279;
case 0x927c:goto L_1927C;
case 0x927e:goto L_1927E;
case 0x927f:goto L_1927F;
case 0x9282:goto L_19282;
case 0x9284:goto L_19284;
case 0x9289:goto L_19289;
case 0x928a:goto L_1928A;
case 0x928d:goto L_1928D;
case 0x9292:goto L_19292;
case 0x9293:goto L_19293;
case 0x9294:goto L_19294;
case 0x9299:goto L_19299;
case 0x929c:goto L_1929C;
case 0x929e:goto L_1929E;
case 0x92a0:goto L_192A0;
case 0x92a3:goto L_192A3;
case 0x92a5:goto L_192A5;
case 0x92a6:goto L_192A6;
case 0x92a7:goto L_192A7;
case 0x92a8:goto L_192A8;
case 0x92a9:goto L_192A9;
case 0x92ab:goto L_192AB;
case 0x92ae:goto L_192AE;
case 0x92b0:goto L_192B0;
case 0x92b3:goto L_192B3;
case 0x92b6:goto L_192B6;
case 0x92b8:goto L_192B8;
case 0x92bb:goto L_192BB;
case 0x92be:goto L_192BE;
case 0x92bf:goto L_192BF;
case 0x92c0:goto L_192C0;
case 0x92c1:goto L_192C1;
case 0x92c4:goto L_192C4;
case 0x92c7:goto L_192C7;
case 0x92c8:goto L_192C8;
case 0x92cb:goto L_192CB;
case 0x92cd:goto L_192CD;
case 0x92cf:goto L_192CF;
case 0x92d2:goto L_192D2;
case 0x92d4:goto L_192D4;
case 0x92d6:goto L_192D6;
case 0x92d7:goto L_192D7;
case 0x92d8:goto L_192D8;
case 0x92db:goto L_192DB;
case 0x92dd:goto L_192DD;
case 0x92e2:goto L_192E2;
case 0x92e3:goto L_192E3;
case 0x92e5:goto L_192E5;
case 0x92e8:goto L_192E8;
case 0x92ea:goto L_192EA;
case 0x92ed:goto L_192ED;
case 0x92ee:goto L_192EE;
case 0x92f0:goto L_192F0;
case 0x92f2:goto L_192F2;
case 0x92f4:goto L_192F4;
case 0x92f5:goto L_192F5;
case 0x92f8:goto L_192F8;
case 0x92fb:goto L_192FB;
case 0x92fe:goto L_192FE;
case 0x9300:goto L_19300;
case 0x9303:goto L_19303;
case 0x9305:goto L_19305;
case 0x9306:goto L_19306;
case 0x9307:goto L_19307;
case 0x9308:goto L_19308;
case 0xeb6c:goto L_1EB6C;
case 0xeb6f:goto L_1EB6F;
case 0xeb71:goto L_1EB71;
case 0xeb74:goto L_1EB74;
case 0xeb75:goto L_1EB75;
case 0xeb78:goto L_1EB78;
case 0xeb7a:goto L_1EB7A;
case 0xeb7b:goto L_1EB7B;
case 0xeb7d:goto L_1EB7D;
case 0xeb7e:goto L_1EB7E;
case 0xeb81:goto L_1EB81;
case 0xeb83:goto L_1EB83;
case 0xeb84:goto L_1EB84;
case 0xeb85:goto L_1EB85;
case 0xeb87:goto L_1EB87;
case 0xeb89:goto L_1EB89;
case 0xeb8b:goto L_1EB8B;
case 0xeb8c:goto L_1EB8C;
case 0xeb8e:goto L_1EB8E;
case 0xeb8f:goto L_1EB8F;
case 0xeb91:goto L_1EB91;
case 0xeb92:goto L_1EB92;
case 0xeb93:goto L_1EB93;
case 0xeb94:goto L_1EB94;
case 0xeb95:goto L_1EB95;
case 0xeb96:goto L_1EB96;
case 0xeb99:goto L_1EB99;
case 0xeb9b:goto L_1EB9B;
case 0xeb9e:goto L_1EB9E;
case 0xeb9f:goto L_1EB9F;
case 0xeba2:goto L_1EBA2;
case 0xeba4:goto L_1EBA4;
case 0xeba5:goto L_1EBA5;
case 0xeba7:goto L_1EBA7;
case 0xeba8:goto L_1EBA8;
case 0xebaa:goto L_1EBAA;
case 0xebab:goto L_1EBAB;
case 0xebae:goto L_1EBAE;
case 0xebb0:goto L_1EBB0;
case 0xebb1:goto L_1EBB1;
case 0xebb2:goto L_1EBB2;
case 0xebb4:goto L_1EBB4;
case 0xebb5:goto L_1EBB5;
case 0xebb6:goto L_1EBB6;
case 0xebb8:goto L_1EBB8;
case 0xebb9:goto L_1EBB9;
case 0xebba:goto L_1EBBA;
case 0xebbc:goto L_1EBBC;
case 0xebbd:goto L_1EBBD;
case 0xebbe:goto L_1EBBE;
case 0xebc0:goto L_1EBC0;
case 0xebc1:goto L_1EBC1;
case 0xebc2:goto L_1EBC2;
case 0xebc4:goto L_1EBC4;
case 0xebc5:goto L_1EBC5;
case 0xebc6:goto L_1EBC6;
case 0xebc8:goto L_1EBC8;
case 0xebc9:goto L_1EBC9;
case 0xebca:goto L_1EBCA;
case 0xebcc:goto L_1EBCC;
case 0xebcd:goto L_1EBCD;
case 0xebcf:goto L_1EBCF;
case 0xebd2:goto L_1EBD2;
case 0xebd3:goto L_1EBD3;
case 0xebd6:goto L_1EBD6;
case 0xebd8:goto L_1EBD8;
case 0xebd9:goto L_1EBD9;
case 0xebda:goto L_1EBDA;
case 0xebdb:goto L_1EBDB;
default:break;}
if(!su_invoke_old)return;
ki_army_invoke(m,m->ip,h);su_invoke_old=0;goto SU_DISPATCH;
L_10000: ; /* IDA 0x10000; mov     ah, 51h */
m->ip=0x0;
/* Original entry: start */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x0,h->user);
vg_ah(m,(uint8_t)(81));
goto L_10002;
L_10002: ; /* IDA 0x10002; int     21h; DOS - 2+ internal - GET PSP SEGMENT */
m->ip=0x2;
su_skip_first_enter=0;
su_interrupt(m,33,0x4,h);
goto L_10004;
L_10004: ; /* IDA 0x10004; mov     ds, bx */
m->ip=0x4;
su_skip_first_enter=0;
m->ds=(uint16_t)(m->bx);
goto L_10006;
L_10006: ; /* IDA 0x10006; mov     al, ds:80h */
m->ip=0x6;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,128)));
goto L_10009;
L_10009: ; /* IDA 0x10009; mov     cs:byte_10D36, al */
m->ip=0x9;
su_skip_first_enter=0;
dl_write8(m,m->cs,3382,(uint8_t)m->ax);
goto L_1000D;
L_1000D: ; /* IDA 0x1000D; mov     ax, cs */
m->ip=0xd;
su_skip_first_enter=0;
m->ax=(uint16_t)(m->cs);
goto L_1000F;
L_1000F: ; /* IDA 0x1000F; mov     ds, ax */
m->ip=0xf;
su_skip_first_enter=0;
m->ds=(uint16_t)(m->ax);
goto L_10011;
L_10011: ; /* IDA 0x10011; mov     es, ax */
m->ip=0x11;
su_skip_first_enter=0;
m->es=(uint16_t)(m->ax);
goto L_10013;
L_10013: ; /* IDA 0x10013; call    sub_100DF */
m->ip=0x13;
su_skip_first_enter=0;
{uint16_t target=223;ec_push(m,0x16);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10016: ; /* IDA 0x10016; mov     ax, 0 */
m->ip=0x16;
su_skip_first_enter=0;
m->ax=(uint16_t)(0);
goto L_10019;
L_10019: ; /* IDA 0x10019; call    sub_20000 */
m->ip=0x19;
su_skip_first_enter=0;
mx_far(m,0x19,0x1e,h);su_invoke_old=0;goto SU_DISPATCH;
L_1001E: ; /* IDA 0x1001E; and     ax, ax */
m->ip=0x1e;
su_skip_first_enter=0;
{uint16_t result=((m->ax)&(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
goto L_10020;
L_10020: ; /* IDA 0x10020; jnz     short loc_10031 */
m->ip=0x20;
su_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_10031;}
goto L_10022;
L_10022: ; /* IDA 0x10022; push    ax */
m->ip=0x22;
su_skip_first_enter=0;
ec_push(m,m->ax);
goto L_10023;
L_10023: ; /* IDA 0x10023; push    dx */
m->ip=0x23;
su_skip_first_enter=0;
ec_push(m,m->dx);
goto L_10024;
L_10024: ; /* IDA 0x10024; mov     dx, 19A6h */
m->ip=0x24;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 2
m->dx=(uint16_t)(6567);
#else
m->dx=(uint16_t)(6566);
#endif
goto L_10027;
L_10027: ; /* IDA 0x10027; mov     ah, 9 */
m->ip=0x27;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(9));
goto L_10029;
L_10029: ; /* IDA 0x10029; int     21h; DOS - PRINT STRING */
m->ip=0x29;
su_skip_first_enter=0;
su_interrupt(m,33,0x2b,h);
goto L_1002B;
L_1002B: ; /* IDA 0x1002B; pop     dx */
m->ip=0x2b;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_pop(m));
goto L_1002C;
L_1002C: ; /* IDA 0x1002C; pop     ax */
m->ip=0x2c;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_1002D;
L_1002D: ; /* IDA 0x1002D; mov     ah, 4Ch */
m->ip=0x2d;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 1
vg_ah(m,(uint8_t)(76));vg_al(m,(uint8_t)(1));
#else
vg_ah(m,(uint8_t)(76));
#endif
goto L_1002F;
L_1002F: ; /* IDA 0x1002F; int     21h; DOS - 2+ - QUIT WITH EXIT CODE (EXIT) */
m->ip=0x2f;
su_skip_first_enter=0;
su_interrupt(m,33,0x31,h);
goto L_10031;
L_10031: ; /* IDA 0x10031; call    sub_1033B */
m->ip=0x31;
su_skip_first_enter=0;
{uint16_t target=827;ec_push(m,0x34);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10034: ; /* IDA 0x10034; call    sub_10210 */
m->ip=0x34;
su_skip_first_enter=0;
{uint16_t target=528;ec_push(m,0x37);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10037: ; /* IDA 0x10037; call    sub_10155 */
m->ip=0x37;
su_skip_first_enter=0;
{uint16_t target=341;ec_push(m,0x3a);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1003A: ; /* IDA 0x1003A; call    sub_1006B */
m->ip=0x3a;
su_skip_first_enter=0;
{uint16_t target=107;ec_push(m,0x3d);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1003D: ; /* IDA 0x1003D; call    sub_1030F */
m->ip=0x3d;
su_skip_first_enter=0;
{uint16_t target=783;ec_push(m,0x40);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10040: ; /* IDA 0x10040; call    sub_1005B */
m->ip=0x40;
su_skip_first_enter=0;
{uint16_t target=91;ec_push(m,0x43);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10043: ; /* IDA 0x10043; call    sub_102C2 */
m->ip=0x43;
/* Original entry: raw_entry_10043 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x43,h->user);
{uint16_t target=706;ec_push(m,0x46);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10046: ; /* IDA 0x10046; push    ax */
m->ip=0x46;
su_skip_first_enter=0;
ec_push(m,m->ax);
goto L_10047;
L_10047: ; /* IDA 0x10047; xor     ax, ax */
m->ip=0x47;
su_skip_first_enter=0;
{uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
goto L_10049;
L_10049: ; /* IDA 0x10049; int     33h; - MS MOUSE - RESET DRIVER AND READ STATUS */
m->ip=0x49;
su_skip_first_enter=0;
su_interrupt(m,51,0x4b,h);
goto L_1004B;
L_1004B: ; /* IDA 0x1004B; call    sub_10236 */
m->ip=0x4b;
su_skip_first_enter=0;
{uint16_t target=566;ec_push(m,0x4e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1004E: ; /* IDA 0x1004E; call    sub_1034E */
m->ip=0x4e;
su_skip_first_enter=0;
{uint16_t target=846;ec_push(m,0x51);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10051: ; /* IDA 0x10051; mov     ax, 3 */
m->ip=0x51;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 12
m->ax=(uint16_t)(2);
#else
m->ax=(uint16_t)(3);
#endif
goto L_10054;
L_10054: ; /* IDA 0x10054; int     10h; - VIDEO - SET VIDEO MODE */
m->ip=0x54;
su_skip_first_enter=0;
su_interrupt(m,16,0x56,h);
goto L_10056;
L_10056: ; /* IDA 0x10056; pop     ax */
m->ip=0x56;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_10057;
L_10057: ; /* IDA 0x10057; mov     ah, 4Ch */
m->ip=0x57;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(76));
goto L_10059;
L_10059: ; /* IDA 0x10059; int     21h; DOS - 2+ - QUIT WITH EXIT CODE (EXIT) */
m->ip=0x59;
su_skip_first_enter=0;
su_interrupt(m,33,0x5b,h);
goto L_1005B;
L_1005B: ; /* IDA 0x1005B; call    sub_10A1C */
m->ip=0x5b;
/* Original entry: sub_1005B */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x5b,h->user);
{uint16_t target=2588;ec_push(m,0x5e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1005E: ; /* IDA 0x1005E; call    sub_119CA */
m->ip=0x5e;
su_skip_first_enter=0;
{uint16_t target=6602;ec_push(m,0x61);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10061: ; /* IDA 0x10061; call    sub_11A6E */
m->ip=0x61;
su_skip_first_enter=0;
{uint16_t target=6766;ec_push(m,0x64);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10064: ; /* IDA 0x10064; call    sub_119CA */
m->ip=0x64;
su_skip_first_enter=0;
{uint16_t target=6602;ec_push(m,0x67);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10067: ; /* IDA 0x10067; call    sub_11BE0 */
m->ip=0x67;
su_skip_first_enter=0;
{uint16_t target=7136;ec_push(m,0x6a);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1006A: ; /* IDA 0x1006A; retn */
m->ip=0x6a;
/* Original entry: raw_entry_1006A */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x6a,h->user);
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1006B: ; /* IDA 0x1006B; cld */
m->ip=0x6b;
/* Original entry: sub_1006B */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x6b,h->user);
m->flags&=(uint16_t)~0x400u;
goto L_1006C;
L_1006C: ; /* IDA 0x1006C; call    sub_1EB6C */
m->ip=0x6c;
su_skip_first_enter=0;
{uint16_t target=60268;ec_push(m,0x6f);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1006F: ; /* IDA 0x1006F; mov     ah, 0 */
m->ip=0x6f;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(0));
goto L_10071;
L_10071: ; /* IDA 0x10071; call    sub_1EB93 */
m->ip=0x71;
su_skip_first_enter=0;
{uint16_t target=60307;ec_push(m,0x74);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10074: ; /* IDA 0x10074; call    sub_10A1C */
m->ip=0x74;
su_skip_first_enter=0;
{uint16_t target=2588;ec_push(m,0x77);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10077: ; /* IDA 0x10077; call    sub_1EC82 */
m->ip=0x77;
su_skip_first_enter=0;
{uint16_t target=60546;ec_push(m,0x7a);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1007A: ; /* IDA 0x1007A; call    sub_1F720 */
m->ip=0x7a;
su_skip_first_enter=0;
{uint16_t target=63264;ec_push(m,0x7d);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1007D: ; /* IDA 0x1007D; mov     al, 0 */
m->ip=0x7d;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(0));
goto L_1007F;
L_1007F: ; /* IDA 0x1007F; call    sub_109AF */
m->ip=0x7f;
su_skip_first_enter=0;
{uint16_t target=2479;ec_push(m,0x82);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_10082: ; /* IDA 0x10082; mov     ax, ds */
m->ip=0x82;
su_skip_first_enter=0;
m->ax=(uint16_t)(m->ds);
goto L_10084;
L_10084: ; /* IDA 0x10084; mov     es, ax */
m->ip=0x84;
su_skip_first_enter=0;
m->es=(uint16_t)(m->ax);
goto L_10086;
L_10086: ; /* IDA 0x10086; mov     si, 18D4h */
m->ip=0x86;
su_skip_first_enter=0;
m->si=(uint16_t)(6356);
goto L_10089;
L_10089: ; /* IDA 0x10089; mov     di, 1964h */
m->ip=0x89;
su_skip_first_enter=0;
m->di=(uint16_t)(6500);
goto L_1008C;
L_1008C: ; /* IDA 0x1008C; mov     cx, 18h */
m->ip=0x8c;
su_skip_first_enter=0;
m->cx=(uint16_t)(24);
goto L_1008F;
L_1008F: ; /* IDA 0x1008F; cld */
m->ip=0x8f;
su_skip_first_enter=0;
m->flags&=(uint16_t)~0x400u;
goto L_10090;
L_10090: ; /* IDA 0x10090; rep movsw */
m->ip=0x90;
su_skip_first_enter=0;
mc_string(m,0,1);
goto L_10092;
L_10092: ; /* IDA 0x10092; mov     dx, 0D64h */
m->ip=0x92;
su_skip_first_enter=0;
m->dx=(uint16_t)(3428);
goto L_10095;
L_10095: ; /* IDA 0x10095; mov     bx, word_10D38 */
m->ip=0x95;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,3384));
goto L_10099;
L_10099: ; /* IDA 0x10099; xor     si, si */
m->ip=0x99;
su_skip_first_enter=0;
{uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_1009B;
L_1009B: ; /* IDA 0x1009B; call    sub_1E378 */
m->ip=0x9b;
su_skip_first_enter=0;
{uint16_t target=58232;ec_push(m,0x9e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1009E: ; /* IDA 0x1009E; mov     dx, 0D6Dh */
m->ip=0x9e;
su_skip_first_enter=0;
m->dx=(uint16_t)(3437);
goto L_100A1;
L_100A1: ; /* IDA 0x100A1; mov     bx, word_10D3A */
m->ip=0xa1;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,3386));
goto L_100A5;
L_100A5: ; /* IDA 0x100A5; xor     si, si */
m->ip=0xa5;
su_skip_first_enter=0;
{uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_100A7;
L_100A7: ; /* IDA 0x100A7; xor     cx, cx */
m->ip=0xa7;
su_skip_first_enter=0;
{uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
goto L_100A9;
L_100A9: ; /* IDA 0x100A9; mov     ax, 6700h */
m->ip=0xa9;
su_skip_first_enter=0;
m->ax=(uint16_t)(26368);
goto L_100AC;
L_100AC: ; /* IDA 0x100AC; mov     di, 3000h */
m->ip=0xac;
su_skip_first_enter=0;
m->di=(uint16_t)(12288);
goto L_100AF;
L_100AF: ; /* IDA 0x100AF; call    sub_1E38C */
m->ip=0xaf;
su_skip_first_enter=0;
{uint16_t target=58252;ec_push(m,0xb2);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_100B2: ; /* IDA 0x100B2; mov     dx, 0D6Dh */
m->ip=0xb2;
su_skip_first_enter=0;
m->dx=(uint16_t)(3437);
goto L_100B5;
L_100B5: ; /* IDA 0x100B5; mov     bx, word_10D48 */
m->ip=0xb5;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,3400));
goto L_100B9;
L_100B9: ; /* IDA 0x100B9; xor     si, si */
m->ip=0xb9;
su_skip_first_enter=0;
{uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_100BB;
L_100BB: ; /* IDA 0x100BB; xor     cx, cx */
m->ip=0xbb;
su_skip_first_enter=0;
{uint16_t result=((m->cx)^(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
goto L_100BD;
L_100BD: ; /* IDA 0x100BD; mov     ax, 9700h */
m->ip=0xbd;
su_skip_first_enter=0;
m->ax=(uint16_t)(38656);
goto L_100C0;
L_100C0: ; /* IDA 0x100C0; mov     di, 0FFFFh */
m->ip=0xc0;
su_skip_first_enter=0;
m->di=(uint16_t)(65535);
goto L_100C3;
L_100C3: ; /* IDA 0x100C3; call    sub_1E38C */
m->ip=0xc3;
su_skip_first_enter=0;
{uint16_t target=58252;ec_push(m,0xc6);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_100C6: ; /* IDA 0x100C6; mov     dx, 0DE4h */
m->ip=0xc6;
su_skip_first_enter=0;
m->dx=(uint16_t)(3556);
goto L_100C9;
L_100C9: ; /* IDA 0x100C9; mov     bx, word_10D3E */
m->ip=0xc9;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,3390));
goto L_100CD;
L_100CD: ; /* IDA 0x100CD; xor     si, si */
m->ip=0xcd;
su_skip_first_enter=0;
{uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_100CF;
L_100CF: ; /* IDA 0x100CF; call    sub_1E378 */
m->ip=0xcf;
su_skip_first_enter=0;
{uint16_t target=58232;ec_push(m,0xd2);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_100D2: ; /* IDA 0x100D2; mov     dx, 0DAFh */
m->ip=0xd2;
su_skip_first_enter=0;
m->dx=(uint16_t)(3503);
goto L_100D5;
L_100D5: ; /* IDA 0x100D5; mov     bx, word_10D44 */
m->ip=0xd5;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,3396));
goto L_100D9;
L_100D9: ; /* IDA 0x100D9; xor     si, si */
m->ip=0xd9;
su_skip_first_enter=0;
{uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_100DB;
L_100DB: ; /* IDA 0x100DB; call    sub_1E364 */
m->ip=0xdb;
su_skip_first_enter=0;
{uint16_t target=58212;ec_push(m,0xde);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_100DE: ; /* IDA 0x100DE; retn */
m->ip=0xde;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_10155: ; /* IDA 0x10155; mov     ax, 0Fh */
m->ip=0x155;
/* Original entry: sub_10155 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x155,h->user);
m->ax=(uint16_t)(15);
goto L_10158;
L_10158: ; /* IDA 0x10158; call    sub_20000 */
m->ip=0x158;
su_skip_first_enter=0;
mx_far(m,0x158,0x15d,h);su_invoke_old=0;goto SU_DISPATCH;
L_1015D: ; /* IDA 0x1015D; mov     ax, 0Dh */
m->ip=0x15d;
su_skip_first_enter=0;
m->ax=(uint16_t)(13);
goto L_10160;
L_10160: ; /* IDA 0x10160; xor     dx, dx */
m->ip=0x160;
su_skip_first_enter=0;
{uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
goto L_10162;
L_10162: ; /* IDA 0x10162; call    sub_20000 */
m->ip=0x162;
su_skip_first_enter=0;
mx_far(m,0x162,0x167,h);su_invoke_old=0;goto SU_DISPATCH;
L_10167: ; /* IDA 0x10167; mov     ax, 0Eh */
m->ip=0x167;
su_skip_first_enter=0;
m->ax=(uint16_t)(14);
goto L_1016A;
L_1016A: ; /* IDA 0x1016A; mov     cx, 1 */
m->ip=0x16a;
su_skip_first_enter=0;
m->cx=(uint16_t)(1);
goto L_1016D;
L_1016D: ; /* IDA 0x1016D; mov     dx, cx */
m->ip=0x16d;
su_skip_first_enter=0;
m->dx=(uint16_t)(m->cx);
goto L_1016F;
L_1016F: ; /* IDA 0x1016F; call    sub_20000 */
m->ip=0x16f;
su_skip_first_enter=0;
mx_far(m,0x16f,0x174,h);su_invoke_old=0;goto SU_DISPATCH;
L_10174: ; /* IDA 0x10174; mov     ax, 5 */
m->ip=0x174;
su_skip_first_enter=0;
m->ax=(uint16_t)(5);
goto L_10177;
L_10177: ; /* IDA 0x10177; xor     bx, bx */
m->ip=0x177;
su_skip_first_enter=0;
{uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
goto L_10179;
L_10179: ; /* IDA 0x10179; call    sub_20000 */
m->ip=0x179;
su_skip_first_enter=0;
mx_far(m,0x179,0x17e,h);su_invoke_old=0;goto SU_DISPATCH;
L_1017E: ; /* IDA 0x1017E; mov     ax, 5 */
m->ip=0x17e;
su_skip_first_enter=0;
m->ax=(uint16_t)(5);
goto L_10181;
L_10181: ; /* IDA 0x10181; mov     bx, 1 */
m->ip=0x181;
su_skip_first_enter=0;
m->bx=(uint16_t)(1);
goto L_10184;
L_10184: ; /* IDA 0x10184; call    sub_20000 */
m->ip=0x184;
su_skip_first_enter=0;
mx_far(m,0x184,0x189,h);su_invoke_old=0;goto SU_DISPATCH;
L_10189: ; /* IDA 0x10189; mov     ax, 7 */
m->ip=0x189;
su_skip_first_enter=0;
m->ax=(uint16_t)(7);
goto L_1018C;
L_1018C: ; /* IDA 0x1018C; mov     cx, 0 */
m->ip=0x18c;
su_skip_first_enter=0;
m->cx=(uint16_t)(0);
goto L_1018F;
L_1018F: ; /* IDA 0x1018F; mov     dx, 27Fh */
m->ip=0x18f;
su_skip_first_enter=0;
m->dx=(uint16_t)(639);
goto L_10192;
L_10192: ; /* IDA 0x10192; call    sub_20000 */
m->ip=0x192;
su_skip_first_enter=0;
mx_far(m,0x192,0x197,h);su_invoke_old=0;goto SU_DISPATCH;
L_10197: ; /* IDA 0x10197; mov     ax, 8 */
m->ip=0x197;
su_skip_first_enter=0;
m->ax=(uint16_t)(8);
goto L_1019A;
L_1019A: ; /* IDA 0x1019A; mov     cx, 0 */
m->ip=0x19a;
su_skip_first_enter=0;
m->cx=(uint16_t)(0);
goto L_1019D;
L_1019D: ; /* IDA 0x1019D; mov     dx, 18Fh */
m->ip=0x19d;
su_skip_first_enter=0;
m->dx=(uint16_t)(399);
goto L_101A0;
L_101A0: ; /* IDA 0x101A0; call    sub_20000 */
m->ip=0x1a0;
su_skip_first_enter=0;
mx_far(m,0x1a0,0x1a5,h);su_invoke_old=0;goto SU_DISPATCH;
L_101A5: ; /* IDA 0x101A5; mov     ax, 4 */
m->ip=0x1a5;
su_skip_first_enter=0;
m->ax=(uint16_t)(4);
goto L_101A8;
L_101A8: ; /* IDA 0x101A8; mov     cx, 140h */
m->ip=0x1a8;
su_skip_first_enter=0;
m->cx=(uint16_t)(320);
goto L_101AB;
L_101AB: ; /* IDA 0x101AB; mov     dx, 0C8h */
m->ip=0x1ab;
su_skip_first_enter=0;
m->dx=(uint16_t)(200);
goto L_101AE;
L_101AE: ; /* IDA 0x101AE; call    sub_20000 */
m->ip=0x1ae;
su_skip_first_enter=0;
mx_far(m,0x1ae,0x1b3,h);su_invoke_old=0;goto SU_DISPATCH;
L_101B3: ; /* IDA 0x101B3; retn */
m->ip=0x1b3;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_10210: ; /* IDA 0x10210; push    ds */
m->ip=0x210;
/* Original entry: sub_10210 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x210,h->user);
ec_push(m,m->ds);
goto L_10211;
L_10211: ; /* IDA 0x10211; push    ax */
m->ip=0x211;
su_skip_first_enter=0;
ec_push(m,m->ax);
goto L_10212;
L_10212: ; /* IDA 0x10212; push    si */
m->ip=0x212;
su_skip_first_enter=0;
ec_push(m,m->si);
goto L_10213;
L_10213: ; /* IDA 0x10213; xor     ah, ah */
m->ip=0x213;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
goto L_10215;
L_10215: ; /* IDA 0x10215; int     61h; reserved for user interrupt */
m->ip=0x215;
su_skip_first_enter=0;
su_interrupt(m,97,0x217,h);
goto L_10217;
L_10217: ; /* IDA 0x10217; mov     ah, 2 */
m->ip=0x217;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 3
vg_ah(m,(uint8_t)(3));
#else
vg_ah(m,(uint8_t)(2));
#endif
goto L_10219;
L_10219: ; /* IDA 0x10219; mov     ds, cs:word_10D3E */
m->ip=0x219;
su_skip_first_enter=0;
m->ds=(uint16_t)(dl_read16(m,m->cs,3390));
goto L_1021E;
L_1021E: ; /* IDA 0x1021E; xor     si, si */
m->ip=0x21e;
su_skip_first_enter=0;
{uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_10220;
L_10220: ; /* IDA 0x10220; int     61h; reserved for user interrupt */
m->ip=0x220;
su_skip_first_enter=0;
su_interrupt(m,97,0x222,h);
goto L_10222;
L_10222: ; /* IDA 0x10222; mov     ah, 3 */
m->ip=0x222;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(3));
goto L_10224;
L_10224: ; /* IDA 0x10224; int     61h; reserved for user interrupt */
m->ip=0x224;
su_skip_first_enter=0;
su_interrupt(m,97,0x226,h);
goto L_10226;
L_10226: ; /* IDA 0x10226; mov     cs:byte_1020E, 0FFh */
m->ip=0x226;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 13
dl_write8(m,m->cs,526,254);
#else
dl_write8(m,m->cs,526,255);
#endif
goto L_1022C;
L_1022C: ; /* IDA 0x1022C; mov     cs:byte_1020F, 0FFh */
m->ip=0x22c;
su_skip_first_enter=0;
dl_write8(m,m->cs,527,255);
goto L_10232;
L_10232: ; /* IDA 0x10232; pop     si */
m->ip=0x232;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_pop(m));
goto L_10233;
L_10233: ; /* IDA 0x10233; pop     ax */
m->ip=0x233;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_10234;
L_10234: ; /* IDA 0x10234; pop     ds */
m->ip=0x234;
su_skip_first_enter=0;
m->ds=(uint16_t)(ec_pop(m));
goto L_10235;
L_10235: ; /* IDA 0x10235; retn */
m->ip=0x235;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_10236: ; /* IDA 0x10236; push    ax */
m->ip=0x236;
/* Original entry: sub_10236 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x236,h->user);
ec_push(m,m->ax);
goto L_10237;
L_10237: ; /* IDA 0x10237; mov     ah, 8 */
m->ip=0x237;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(8));
goto L_10239;
L_10239: ; /* IDA 0x10239; int     61h; reserved for user interrupt */
m->ip=0x239;
su_skip_first_enter=0;
su_interrupt(m,97,0x23b,h);
goto L_1023B;
L_1023B: ; /* IDA 0x1023B; mov     ah, 4 */
m->ip=0x23b;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(4));
goto L_1023D;
L_1023D: ; /* IDA 0x1023D; int     61h; reserved for user interrupt */
m->ip=0x23d;
su_skip_first_enter=0;
su_interrupt(m,97,0x23f,h);
goto L_1023F;
L_1023F: ; /* IDA 0x1023F; pop     ax */
m->ip=0x23f;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_10240;
L_10240: ; /* IDA 0x10240; retn */
m->ip=0x240;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1033B: ; /* IDA 0x1033B; push    ds */
m->ip=0x33b;
/* Original entry: sub_1033B */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x33b,h->user);
ec_push(m,m->ds);
goto L_1033C;
L_1033C: ; /* IDA 0x1033C; push    ax */
m->ip=0x33c;
su_skip_first_enter=0;
ec_push(m,m->ax);
goto L_1033D;
L_1033D: ; /* IDA 0x1033D; push    dx */
m->ip=0x33d;
su_skip_first_enter=0;
ec_push(m,m->dx);
goto L_1033E;
L_1033E: ; /* IDA 0x1033E; mov     ax, cs */
m->ip=0x33e;
su_skip_first_enter=0;
m->ax=(uint16_t)(m->cs);
goto L_10340;
L_10340: ; /* IDA 0x10340; mov     ds, ax */
m->ip=0x340;
su_skip_first_enter=0;
m->ds=(uint16_t)(m->ax);
goto L_10342;
L_10342: ; /* IDA 0x10342; mov     dx, 356h */
m->ip=0x342;
su_skip_first_enter=0;
m->dx=(uint16_t)(854);
goto L_10345;
L_10345: ; /* IDA 0x10345; mov     ax, 0C00h */
m->ip=0x345;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 5
m->ax=(uint16_t)(3073);
#else
m->ax=(uint16_t)(3072);
#endif
goto L_10348;
L_10348: ; /* IDA 0x10348; int     61h; reserved for user interrupt */
m->ip=0x348;
su_skip_first_enter=0;
su_interrupt(m,97,0x34a,h);
goto L_1034A;
L_1034A: ; /* IDA 0x1034A; pop     dx */
m->ip=0x34a;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_pop(m));
goto L_1034B;
L_1034B: ; /* IDA 0x1034B; pop     ax */
m->ip=0x34b;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_1034C;
L_1034C: ; /* IDA 0x1034C; pop     ds */
m->ip=0x34c;
su_skip_first_enter=0;
m->ds=(uint16_t)(ec_pop(m));
goto L_1034D;
L_1034D: ; /* IDA 0x1034D; retn */
m->ip=0x34d;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1034E: ; /* IDA 0x1034E; push    ax */
m->ip=0x34e;
/* Original entry: sub_1034E */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x34e,h->user);
ec_push(m,m->ax);
goto L_1034F;
L_1034F: ; /* IDA 0x1034F; mov     ax, 0C01h */
m->ip=0x34f;
su_skip_first_enter=0;
m->ax=(uint16_t)(3073);
goto L_10352;
L_10352: ; /* IDA 0x10352; int     61h; reserved for user interrupt */
m->ip=0x352;
su_skip_first_enter=0;
su_interrupt(m,97,0x354,h);
goto L_10354;
L_10354: ; /* IDA 0x10354; pop     ax */
m->ip=0x354;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_10355;
L_10355: ; /* IDA 0x10355; retn */
m->ip=0x355;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_10356: ; /* IDA 0x10356; cli */
m->ip=0x356;
/* Original entry: raw_entry_10356 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x356,h->user);
m->flags&=(uint16_t)~0x200u;
goto L_10357;
L_10357: ; /* IDA 0x10357; pushf */
m->ip=0x357;
su_skip_first_enter=0;
ec_push(m,m->flags);
goto L_10358;
L_10358: ; /* IDA 0x10358; mov     cs:byte_10D2C, 1 */
m->ip=0x358;
su_skip_first_enter=0;
dl_write8(m,m->cs,3372,1);
goto L_1035E;
L_1035E: ; /* IDA 0x1035E; cmp     cs:byte_10D2D, 0FFh */
m->ip=0x35e;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 9
ec_cmp(m,dl_read8(m,m->cs,3373),254,8);
#else
ec_cmp(m,dl_read8(m,m->cs,3373),255,8);
#endif
goto L_10364;
L_10364: ; /* IDA 0x10364; jnb     short loc_1036B */
m->ip=0x364;
su_skip_first_enter=0;
if(!(m->flags&1)) {goto L_1036B;}
goto L_10366;
L_10366: ; /* IDA 0x10366; inc     cs:byte_10D2D */
m->ip=0x366;
su_skip_first_enter=0;
dl_write8(m,m->cs,3373,ec_inc(m,dl_read8(m,m->cs,3373),8));
goto L_1036B;
L_1036B: ; /* IDA 0x1036B; popf */
m->ip=0x36b;
su_skip_first_enter=0;
m->flags=ec_pop(m)|2;
goto L_1036C;
L_1036C: ; /* IDA 0x1036C; retf */
m->ip=0x36c;
su_skip_first_enter=0;
m->ip=ec_pop(m);m->cs=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_10612: ; /* IDA 0x10612; mov     ds, cx */
m->ip=0x612;
/* Original entry: sub_10612 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x612,h->user);
m->ds=(uint16_t)(m->cx);
goto L_10614;
L_10614: ; /* IDA 0x10614; mov     es, cx */
m->ip=0x614;
su_skip_first_enter=0;
m->es=(uint16_t)(m->cx);
goto L_10616;
L_10616: ; /* IDA 0x10616; mov     cx, ax */
m->ip=0x616;
su_skip_first_enter=0;
m->cx=(uint16_t)(m->ax);
goto L_10618;
L_10618: ; /* IDA 0x10618; mov     si, bx */
m->ip=0x618;
su_skip_first_enter=0;
m->si=(uint16_t)(m->bx);
goto L_1061A;
L_1061A: ; /* IDA 0x1061A; mov     di, dx */
m->ip=0x61a;
su_skip_first_enter=0;
m->di=(uint16_t)(m->dx);
goto L_1061C;
L_1061C: ; /* IDA 0x1061C; rep movsw */
m->ip=0x61c;
su_skip_first_enter=0;
mc_string(m,0,1);
goto L_1061E;
L_1061E: ; /* IDA 0x1061E; retn */
m->ip=0x61e;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_11A6E: ; /* IDA 0x11A6E; mov     ax, cs */
m->ip=0x1a6e;
/* Original entry: sub_11A6E */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x1a6e,h->user);
m->ax=(uint16_t)(m->cs);
goto L_11A70;
L_11A70: ; /* IDA 0x11A70; mov     ds, ax */
m->ip=0x1a70;
su_skip_first_enter=0;
m->ds=(uint16_t)(m->ax);
goto L_11A72;
L_11A72: ; /* IDA 0x11A72; mov     es, ax */
m->ip=0x1a72;
su_skip_first_enter=0;
m->es=(uint16_t)(m->ax);
goto L_11A74;
L_11A74: ; /* IDA 0x11A74; mov     al, 0 */
m->ip=0x1a74;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(0));
goto L_11A76;
L_11A76: ; /* IDA 0x11A76; call    sub_10241 */
m->ip=0x1a76;
su_skip_first_enter=0;
{uint16_t target=577;ec_push(m,0x1a79);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11A79: ; /* IDA 0x11A79; xor     ax, ax */
m->ip=0x1a79;
su_skip_first_enter=0;
{uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
goto L_11A7B;
L_11A7B: ; /* IDA 0x11A7B; call    sub_18CAE */
m->ip=0x1a7b;
su_skip_first_enter=0;
{uint16_t target=36014;ec_push(m,0x1a7e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11A7E: ; /* IDA 0x11A7E; mov     word_19892, 0FFFFh */
m->ip=0x1a7e;
su_skip_first_enter=0;
dl_write16(m,m->ds,39058,65535);
goto L_11A84;
L_11A84: ; /* IDA 0x11A84; mov     dx, 0AAh */
m->ip=0x1a84;
su_skip_first_enter=0;
m->dx=(uint16_t)(170);
goto L_11A87;
L_11A87: ; /* IDA 0x11A87; mov     bx, 62h ; 'b' */
m->ip=0x1a87;
su_skip_first_enter=0;
m->bx=(uint16_t)(98);
goto L_11A8A;
L_11A8A: ; /* IDA 0x11A8A; call    sub_1D615 */
m->ip=0x1a8a;
su_skip_first_enter=0;
{uint16_t target=54805;ec_push(m,0x1a8d);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11A8D: ; /* IDA 0x11A8D; call    sub_1D66A */
m->ip=0x1a8d;
su_skip_first_enter=0;
{uint16_t target=54890;ec_push(m,0x1a90);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11A90: ; /* IDA 0x11A90; call    sub_109D0 */
m->ip=0x1a90;
su_skip_first_enter=0;
{uint16_t target=2512;ec_push(m,0x1a93);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11A93: ; /* IDA 0x11A93; call    sub_11AC3 */
m->ip=0x1a93;
su_skip_first_enter=0;
{uint16_t target=6851;ec_push(m,0x1a96);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11A96: ; /* IDA 0x11A96; call    sub_102C2 */
m->ip=0x1a96;
su_skip_first_enter=0;
{uint16_t target=706;ec_push(m,0x1a99);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11A99: ; /* IDA 0x11A99; call    sub_10A1C */
m->ip=0x1a99;
su_skip_first_enter=0;
{uint16_t target=2588;ec_push(m,0x1a9c);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11A9C: ; /* IDA 0x11A9C; xor     ah, ah */
m->ip=0x1a9c;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
goto L_11A9E;
L_11A9E: ; /* IDA 0x11A9E; call    sub_1EB93 */
m->ip=0x1a9e;
su_skip_first_enter=0;
{uint16_t target=60307;ec_push(m,0x1aa1);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AA1: ; /* IDA 0x11AA1; call    sub_11AAE */
m->ip=0x1aa1;
su_skip_first_enter=0;
{uint16_t target=6830;ec_push(m,0x1aa4);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AA4: ; /* IDA 0x11AA4; call    sub_10A1C */
m->ip=0x1aa4;
su_skip_first_enter=0;
{uint16_t target=2588;ec_push(m,0x1aa7);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AA7: ; /* IDA 0x11AA7; mov     cs:byte_198A6, 0 */
m->ip=0x1aa7;
su_skip_first_enter=0;
dl_write8(m,m->cs,39078,0);
goto L_11AAD;
L_11AAD: ; /* IDA 0x11AAD; retn */
m->ip=0x1aad;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_11AAE: ; /* IDA 0x11AAE; call    sub_16097 */
m->ip=0x1aae;
/* Original entry: sub_11AAE */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x1aae,h->user);
{uint16_t target=24727;ec_push(m,0x1ab1);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AB1: ; /* IDA 0x11AB1; call    sub_19321 */
m->ip=0x1ab1;
su_skip_first_enter=0;
{uint16_t target=37665;ec_push(m,0x1ab4);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AB4: ; /* IDA 0x11AB4; call    sub_160A1 */
m->ip=0x1ab4;
su_skip_first_enter=0;
{uint16_t target=24737;ec_push(m,0x1ab7);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AB7: ; /* IDA 0x11AB7; call    sub_160A5 */
m->ip=0x1ab7;
su_skip_first_enter=0;
{uint16_t target=24741;ec_push(m,0x1aba);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11ABA: ; /* IDA 0x11ABA; xor     ax, ax */
m->ip=0x1aba;
su_skip_first_enter=0;
{uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
goto L_11ABC;
L_11ABC: ; /* IDA 0x11ABC; mov     word_198A8, ax */
m->ip=0x1abc;
su_skip_first_enter=0;
dl_write16(m,m->ds,39080,m->ax);
goto L_11ABF;
L_11ABF: ; /* IDA 0x11ABF; mov     word_198AA, ax */
m->ip=0x1abf;
su_skip_first_enter=0;
dl_write16(m,m->ds,39082,m->ax);
goto L_11AC2;
L_11AC2: ; /* IDA 0x11AC2; retn */
m->ip=0x1ac2;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_11AC3: ; /* IDA 0x11AC3; call    sub_12078 */
m->ip=0x1ac3;
/* Original entry: sub_11AC3 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x1ac3,h->user);
{uint16_t target=8312;ec_push(m,0x1ac6);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AC6: ; /* IDA 0x11AC6; mov     si, 98C8h */
m->ip=0x1ac6;
su_skip_first_enter=0;
m->si=(uint16_t)(39112);
goto L_11AC9;
L_11AC9: ; /* IDA 0x11AC9; mov     dx, 0D0h */
m->ip=0x1ac9;
su_skip_first_enter=0;
m->dx=(uint16_t)(208);
goto L_11ACC;
L_11ACC: ; /* IDA 0x11ACC; mov     bx, 80h */
m->ip=0x1acc;
su_skip_first_enter=0;
m->bx=(uint16_t)(128);
goto L_11ACF;
L_11ACF: ; /* IDA 0x11ACF; mov     di, 28h ; '(' */
m->ip=0x1acf;
su_skip_first_enter=0;
m->di=(uint16_t)(40);
goto L_11AD2;
L_11AD2: ; /* IDA 0x11AD2; call    sub_18DC8 */
m->ip=0x1ad2;
su_skip_first_enter=0;
{uint16_t target=36296;ec_push(m,0x1ad5);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AD5: ; /* IDA 0x11AD5; and     ax, ax */
m->ip=0x1ad5;
su_skip_first_enter=0;
{uint16_t result=((m->ax)&(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
goto L_11AD7;
L_11AD7: ; /* IDA 0x11AD7; jz      short loc_11AE0 */
m->ip=0x1ad7;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_11AE0;}
goto L_11AD9;
L_11AD9: ; /* IDA 0x11AD9; call    sub_18B40 */
m->ip=0x1ad9;
su_skip_first_enter=0;
{uint16_t target=35648;ec_push(m,0x1adc);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11ADC: ; /* IDA 0x11ADC; jb      short loc_11AC6 */
m->ip=0x1adc;
su_skip_first_enter=0;
if(m->flags&1) {goto L_11AC6;}
goto L_11ADE;
L_11ADE: ; /* IDA 0x11ADE; jmp     short loc_11B2C */
m->ip=0x1ade;
su_skip_first_enter=0;
goto L_11B2C;
L_11AE0: ; /* IDA 0x11AE0; call    sub_18B12 */
m->ip=0x1ae0;
su_skip_first_enter=0;
{uint16_t target=35602;ec_push(m,0x1ae3);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AE3: ; /* IDA 0x11AE3; jb      short loc_11AC6 */
m->ip=0x1ae3;
su_skip_first_enter=0;
if(m->flags&1) {goto L_11AC6;}
goto L_11AE5;
L_11AE5: ; /* IDA 0x11AE5; call    sub_17B3C */
m->ip=0x1ae5;
su_skip_first_enter=0;
{uint16_t target=31548;ec_push(m,0x1ae8);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AE8: ; /* IDA 0x11AE8; jb      short loc_11AE0 */
m->ip=0x1ae8;
su_skip_first_enter=0;
if(m->flags&1) {goto L_11AE0;}
goto L_11AEA;
L_11AEA: ; /* IDA 0x11AEA; call    sub_18E5A */
m->ip=0x1aea;
su_skip_first_enter=0;
{uint16_t target=36442;ec_push(m,0x1aed);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AED: ; /* IDA 0x11AED; jb      short loc_11AE5 */
m->ip=0x1aed;
su_skip_first_enter=0;
if(m->flags&1) {goto L_11AE5;}
goto L_11AEF;
L_11AEF: ; /* IDA 0x11AEF; and     al, al */
m->ip=0x1aef;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_11AF1;
L_11AF1: ; /* IDA 0x11AF1; jz      short loc_11AF8 */
m->ip=0x1af1;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_11AF8;}
goto L_11AF3;
L_11AF3: ; /* IDA 0x11AF3; call    sub_18FC9 */
m->ip=0x1af3;
su_skip_first_enter=0;
{uint16_t target=36809;ec_push(m,0x1af6);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11AF6: ; /* IDA 0x11AF6; jmp     short loc_11AEA */
m->ip=0x1af6;
su_skip_first_enter=0;
goto L_11AEA;
L_11AF8: ; /* IDA 0x11AF8; mov     ds, cs:word_10D52 */
m->ip=0x1af8;
su_skip_first_enter=0;
m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
goto L_11AFD;
L_11AFD: ; /* IDA 0x11AFD; mov     ah, [bx+2] */
m->ip=0x1afd;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+2))));
goto L_11B00;
L_11B00: ; /* IDA 0x11B00; cmp     ah, 7Fh */
m->ip=0x1b00;
su_skip_first_enter=0;
ec_cmp(m,(uint8_t)(m->ax>>8),127,8);
goto L_11B03;
L_11B03: ; /* IDA 0x11B03; jz      short loc_11B08 */
m->ip=0x1b03;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_11B08;}
goto L_11B05;
L_11B05: ; /* IDA 0x11B05; dec     byte ptr [bx+18h] */
m->ip=0x1b05;
su_skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->bx+24),st_dec(m,dl_read8(m,m->ds,(uint16_t)(m->bx+24)),8));
goto L_11B08;
L_11B08: ; /* IDA 0x11B08; xor     al, al */
m->ip=0x1b08;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_11B0A;
L_11B0A: ; /* IDA 0x11B0A; shr     ax, 1 */
m->ip=0x1b0a;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
goto L_11B0C;
L_11B0C: ; /* IDA 0x11B0C; shr     ax, 1 */
m->ip=0x1b0c;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
goto L_11B0E;
L_11B0E: ; /* IDA 0x11B0E; shr     ax, 1 */
m->ip=0x1b0e;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
goto L_11B10;
L_11B10: ; /* IDA 0x11B10; mov     si, ax */
m->ip=0x1b10;
su_skip_first_enter=0;
m->si=(uint16_t)(m->ax);
goto L_11B12;
L_11B12: ; /* IDA 0x11B12; mov     byte ptr [si+4240h], 0 */
m->ip=0x1b12;
su_skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+16960),0);
goto L_11B17;
L_11B17: ; /* IDA 0x11B17; mov     cs:word_10CFD, bx */
m->ip=0x1b17;
su_skip_first_enter=0;
dl_write16(m,m->cs,3325,m->bx);
goto L_11B1C;
L_11B1C: ; /* IDA 0x11B1C; shl     bx, 1 */
m->ip=0x1b1c;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_11B1E;
L_11B1E: ; /* IDA 0x11B1E; shl     bx, 1 */
m->ip=0x1b1e;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_11B20;
L_11B20: ; /* IDA 0x11B20; mov     cs:byte_10CFF, bh */
m->ip=0x1b20;
su_skip_first_enter=0;
dl_write8(m,m->cs,3327,(uint8_t)(m->bx>>8));
goto L_11B25;
L_11B25: ; /* IDA 0x11B25; mov     ax, cs */
m->ip=0x1b25;
su_skip_first_enter=0;
m->ax=(uint16_t)(m->cs);
goto L_11B27;
L_11B27: ; /* IDA 0x11B27; mov     ds, ax */
m->ip=0x1b27;
su_skip_first_enter=0;
m->ds=(uint16_t)(m->ax);
goto L_11B29;
L_11B29: ; /* IDA 0x11B29; call    sub_12BD9 */
m->ip=0x1b29;
su_skip_first_enter=0;
{uint16_t target=11225;ec_push(m,0x1b2c);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11B2C: ; /* IDA 0x11B2C; call    sub_120D6 */
m->ip=0x1b2c;
su_skip_first_enter=0;
{uint16_t target=8406;ec_push(m,0x1b2f);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11B2F: ; /* IDA 0x11B2F; mov     ds, cs:word_10D52 */
m->ip=0x1b2f;
su_skip_first_enter=0;
m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
goto L_11B34;
L_11B34: ; /* IDA 0x11B34; mov     bx, cs:word_10CFD */
m->ip=0x1b34;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->cs,3325));
goto L_11B39;
L_11B39: ; /* IDA 0x11B39; mov     bh, [bx+3] */
m->ip=0x1b39;
su_skip_first_enter=0;
gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+3))));
goto L_11B3C;
L_11B3C: ; /* IDA 0x11B3C; xor     bl, bl */
m->ip=0x1b3c;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
goto L_11B3E;
L_11B3E: ; /* IDA 0x11B3E; shr     bx, 1 */
m->ip=0x1b3e;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_11B40;
L_11B40: ; /* IDA 0x11B40; shr     bx, 1 */
m->ip=0x1b40;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_11B42;
L_11B42: ; /* IDA 0x11B42; shr     bx, 1 */
m->ip=0x1b42;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_11B44;
L_11B44: ; /* IDA 0x11B44; mov     dx, [bx+848h] */
m->ip=0x1b44;
su_skip_first_enter=0;
m->dx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+2120)));
goto L_11B48;
L_11B48: ; /* IDA 0x11B48; mov     bx, [bx+84Ah] */
m->ip=0x1b48;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,(uint16_t)(m->bx+2122)));
goto L_11B4C;
L_11B4C: ; /* IDA 0x11B4C; mov     ax, 14h */
m->ip=0x1b4c;
su_skip_first_enter=0;
m->ax=(uint16_t)(20);
goto L_11B4F;
L_11B4F: ; /* IDA 0x11B4F; mov     cx, 0Ch */
m->ip=0x1b4f;
su_skip_first_enter=0;
m->cx=(uint16_t)(12);
goto L_11B52;
L_11B52: ; /* IDA 0x11B52; call    sub_12151 */
m->ip=0x1b52;
su_skip_first_enter=0;
{uint16_t target=8529;ec_push(m,0x1b55);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11B55: ; /* IDA 0x11B55; mov     ax, cs */
m->ip=0x1b55;
su_skip_first_enter=0;
m->ax=(uint16_t)(m->cs);
goto L_11B57;
L_11B57: ; /* IDA 0x11B57; mov     ds, ax */
m->ip=0x1b57;
su_skip_first_enter=0;
m->ds=(uint16_t)(m->ax);
goto L_11B59;
L_11B59: ; /* IDA 0x11B59; retn */
m->ip=0x1b59;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_11BE0: ; /* IDA 0x11BE0; mov     ax, cs */
m->ip=0x1be0;
/* Original entry: sub_11BE0 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x1be0,h->user);
m->ax=(uint16_t)(m->cs);
goto L_11BE2;
L_11BE2: ; /* IDA 0x11BE2; mov     ds, ax */
m->ip=0x1be2;
su_skip_first_enter=0;
m->ds=(uint16_t)(m->ax);
goto L_11BE4;
L_11BE4: ; /* IDA 0x11BE4; mov     es, ax */
m->ip=0x1be4;
su_skip_first_enter=0;
m->es=(uint16_t)(m->ax);
goto L_11BE6;
L_11BE6: ; /* IDA 0x11BE6; call    sub_189F0 */
m->ip=0x1be6;
su_skip_first_enter=0;
{uint16_t target=35312;ec_push(m,0x1be9);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11BE9: ; /* IDA 0x11BE9; call    sub_1533D */
m->ip=0x1be9;
su_skip_first_enter=0;
{uint16_t target=21309;ec_push(m,0x1bec);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11BEC: ; /* IDA 0x11BEC; mov     cs:word_19903, ss */
m->ip=0x1bec;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 11
dl_write16(m,m->cs,39171,m->sp);
#else
dl_write16(m,m->cs,39171,m->ss);
#endif
goto L_11BF1;
L_11BF1: ; /* IDA 0x11BF1; mov     cs:word_19901, sp */
m->ip=0x1bf1;
su_skip_first_enter=0;
dl_write16(m,m->cs,39169,m->sp);
goto L_11BF6;
L_11BF6: ; /* IDA 0x11BF6; call    sub_11F7F */
m->ip=0x1bf6;
su_skip_first_enter=0;
{uint16_t target=8063;ec_push(m,0x1bf9);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11BF9: ; /* IDA 0x11BF9; mov     dx, word_1988E */
m->ip=0x1bf9;
su_skip_first_enter=0;
m->dx=(uint16_t)(dl_read16(m,m->ds,39054));
goto L_11BFD;
L_11BFD: ; /* IDA 0x11BFD; mov     bx, word_19890 */
m->ip=0x1bfd;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,39056));
goto L_11C01;
L_11C01: ; /* IDA 0x11C01; call    sub_1D615 */
m->ip=0x1c01;
su_skip_first_enter=0;
{uint16_t target=54805;ec_push(m,0x1c04);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C04: ; /* IDA 0x11C04; call    sub_11CC9 */
m->ip=0x1c04;
su_skip_first_enter=0;
{uint16_t target=7369;ec_push(m,0x1c07);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C07: ; /* IDA 0x11C07; call    sub_1D66A */
m->ip=0x1c07;
su_skip_first_enter=0;
{uint16_t target=54890;ec_push(m,0x1c0a);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C0A: ; /* IDA 0x11C0A; call    sub_109D0 */
m->ip=0x1c0a;
su_skip_first_enter=0;
{uint16_t target=2512;ec_push(m,0x1c0d);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C0D: ; /* IDA 0x11C0D; mov     ax, 2 */
m->ip=0x1c0d;
su_skip_first_enter=0;
m->ax=(uint16_t)(2);
goto L_11C10;
L_11C10: ; /* IDA 0x11C10; mov     byte_198A2, ah */
m->ip=0x1c10;
su_skip_first_enter=0;
dl_write8(m,m->ds,39074,(uint8_t)(m->ax>>8));
goto L_11C14;
L_11C14: ; /* IDA 0x11C14; call    sub_20000 */
m->ip=0x1c14;
su_skip_first_enter=0;
mx_far(m,0x1c14,0x1c19,h);su_invoke_old=0;goto SU_DISPATCH;
L_11C19: ; /* IDA 0x11C19; mov     cs:byte_198A4, 0 */
m->ip=0x1c19;
su_skip_first_enter=0;
dl_write8(m,m->cs,39076,0);
goto L_11C1F;
L_11C1F: ; /* IDA 0x11C1F; call    sub_11F7F */
m->ip=0x1c1f;
/* Original entry: inner pause entry */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x1c1f,h->user);
{uint16_t target=8063;ec_push(m,0x1c22);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C22: ; /* IDA 0x11C22; push    cx */
m->ip=0x1c22;
su_skip_first_enter=0;
ec_push(m,m->cx);
goto L_11C23;
L_11C23: ; /* IDA 0x11C23; push    dx */
m->ip=0x1c23;
su_skip_first_enter=0;
ec_push(m,m->dx);
goto L_11C24;
L_11C24: ; /* IDA 0x11C24; call    sub_1E453 */
m->ip=0x1c24;
su_skip_first_enter=0;
{uint16_t target=58451;ec_push(m,0x1c27);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C27: ; /* IDA 0x11C27; pop     dx */
m->ip=0x1c27;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_pop(m));
goto L_11C28;
L_11C28: ; /* IDA 0x11C28; pop     cx */
m->ip=0x1c28;
su_skip_first_enter=0;
m->cx=(uint16_t)(ec_pop(m));
goto L_11C29;
L_11C29: ; /* IDA 0x11C29; and     al, al */
m->ip=0x1c29;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)&((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_11C2B;
L_11C2B: ; /* IDA 0x11C2B; jz      short loc_11C46 */
m->ip=0x1c2b;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_11C46;}
goto L_11C2D;
L_11C2D: ; /* IDA 0x11C2D; and     bx, bx */
m->ip=0x1c2d;
su_skip_first_enter=0;
{uint16_t result=((m->bx)&(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
goto L_11C2F;
L_11C2F: ; /* IDA 0x11C2F; jz      short loc_11C3E */
m->ip=0x1c2f;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_11C3E;}
goto L_11C31;
L_11C31: ; /* IDA 0x11C31; cmp     bx, 1 */
m->ip=0x1c31;
su_skip_first_enter=0;
ec_cmp(m,m->bx,1,16);
goto L_11C34;
L_11C34: ; /* IDA 0x11C34; jnz     short loc_11C3B */
m->ip=0x1c34;
su_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_11C3B;}
goto L_11C36;
L_11C36: ; /* IDA 0x11C36; call    sub_159A6 */
m->ip=0x1c36;
su_skip_first_enter=0;
{uint16_t target=22950;ec_push(m,0x1c39);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C39: ; /* IDA 0x11C39; jmp     short loc_11C3E */
m->ip=0x1c39;
su_skip_first_enter=0;
goto L_11C3E;
L_11C3B: ; /* IDA 0x11C3B; call    sub_159B7 */
m->ip=0x1c3b;
su_skip_first_enter=0;
{uint16_t target=22967;ec_push(m,0x1c3e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C3E: ; /* IDA 0x11C3E; mov     cs:byte_198A2, 1 */
m->ip=0x1c3e;
su_skip_first_enter=0;
dl_write8(m,m->cs,39074,1);
goto L_11C44;
L_11C44: ; /* IDA 0x11C44; jmp     short loc_11C51 */
m->ip=0x1c44;
su_skip_first_enter=0;
goto L_11C51;
L_11C46: ; /* IDA 0x11C46; cmp     bl, 1 */
m->ip=0x1c46;
su_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->bx,1,8);
goto L_11C49;
L_11C49: ; /* IDA 0x11C49; jnz     short loc_11C4E */
m->ip=0x1c49;
su_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_11C4E;}
goto L_11C4B;
L_11C4B: ; /* IDA 0x11C4B; call    sub_11E46 */
m->ip=0x1c4b;
su_skip_first_enter=0;
{uint16_t target=7750;ec_push(m,0x1c4e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C4E: ; /* IDA 0x11C4E; call    sub_11F30 */
m->ip=0x1c4e;
su_skip_first_enter=0;
{uint16_t target=7984;ec_push(m,0x1c51);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C51: ; /* IDA 0x11C51; test    byte_198A6, 4 */
m->ip=0x1c51;
su_skip_first_enter=0;
{uint16_t result=((dl_read8(m,m->ds,39078))&(4))&255;ec_logic(m,result,8);}
goto L_11C56;
L_11C56: ; /* IDA 0x11C56; jz      short loc_11C5B */
m->ip=0x1c56;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_11C5B;}
goto L_11C58;
L_11C58: ; /* IDA 0x11C58; call    sub_15C58 */
m->ip=0x1c58;
su_skip_first_enter=0;
{uint16_t target=23640;ec_push(m,0x1c5b);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C5B: ; /* IDA 0x11C5B; cmp     byte_198A3, 80h */
m->ip=0x1c5b;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 10
ec_cmp(m,dl_read8(m,m->ds,39075),129,8);
#else
ec_cmp(m,dl_read8(m,m->ds,39075),128,8);
#endif
goto L_11C60;
L_11C60: ; /* IDA 0x11C60; jb      short loc_11C69 */
m->ip=0x1c60;
su_skip_first_enter=0;
if(m->flags&1) {goto L_11C69;}
goto L_11C62;
L_11C62: ; /* IDA 0x11C62; call    sub_11CD0 */
m->ip=0x1c62;
su_skip_first_enter=0;
{uint16_t target=7376;ec_push(m,0x1c65);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C65: ; /* IDA 0x11C65; inc     byte_198A4 */
m->ip=0x1c65;
su_skip_first_enter=0;
dl_write8(m,m->ds,39076,ec_inc(m,dl_read8(m,m->ds,39076),8));
goto L_11C69;
L_11C69: ; /* IDA 0x11C69; mov     dx, word_1988E */
m->ip=0x1c69;
su_skip_first_enter=0;
m->dx=(uint16_t)(dl_read16(m,m->ds,39054));
goto L_11C6D;
L_11C6D: ; /* IDA 0x11C6D; mov     bx, word_19890 */
m->ip=0x1c6d;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,39056));
goto L_11C71;
L_11C71: ; /* IDA 0x11C71; call    sub_1D615 */
m->ip=0x1c71;
su_skip_first_enter=0;
{uint16_t target=54805;ec_push(m,0x1c74);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C74: ; /* IDA 0x11C74; call    sub_11CC9 */
m->ip=0x1c74;
su_skip_first_enter=0;
{uint16_t target=7369;ec_push(m,0x1c77);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C77: ; /* IDA 0x11C77; call    sub_1D66A */
m->ip=0x1c77;
su_skip_first_enter=0;
{uint16_t target=54890;ec_push(m,0x1c7a);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C7A: ; /* IDA 0x11C7A; cmp     byte_198A3, 80h */
m->ip=0x1c7a;
su_skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,39075),128,8);
goto L_11C7F;
L_11C7F: ; /* IDA 0x11C7F; jb      short loc_11C8B */
m->ip=0x1c7f;
su_skip_first_enter=0;
if(m->flags&1) {goto L_11C8B;}
goto L_11C81;
L_11C81: ; /* IDA 0x11C81; cmp     byte_198A2, 0 */
m->ip=0x1c81;
su_skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,39074),0,8);
goto L_11C86;
L_11C86: ; /* IDA 0x11C86; jz      short loc_11C1F */
m->ip=0x1c86;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_11C1F;}
goto L_11C88;
L_11C88: ; /* IDA 0x11C88; call    sub_11F5A */
m->ip=0x1c88;
su_skip_first_enter=0;
{uint16_t target=8026;ec_push(m,0x1c8b);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11C8B: ; /* IDA 0x11C8B; jmp     short loc_11C1F */
m->ip=0x1c8b;
su_skip_first_enter=0;
goto L_11C1F;
L_11CD0: ; /* IDA 0x11CD0; cmp     byte ptr ds:98A5h, 0 */
m->ip=0x1cd0;
/* Original entry: sub_11CD0 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x1cd0,h->user);
ec_cmp(m,dl_read8(m,m->ds,39077),0,8);
goto L_11CD5;
L_11CD5: ; /* IDA 0x11CD5; jz      short loc_11D0B */
m->ip=0x1cd5;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_11D0B;}
goto L_11CD7;
L_11CD7: ; /* IDA 0x11CD7; cmp     byte ptr ds:98A5h, 1 */
m->ip=0x1cd7;
su_skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,39077),1,8);
goto L_11CDC;
L_11CDC: ; /* IDA 0x11CDC; jz      short loc_11CED */
m->ip=0x1cdc;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_11CED;}
goto L_11CDE;
L_11CDE: ; /* IDA 0x11CDE; mov     byte ptr ds:0D2Ch, 0 */
m->ip=0x1cde;
su_skip_first_enter=0;
dl_write8(m,m->ds,3372,0);
goto L_11CE3;
L_11CE3: ; /* IDA 0x11CE3; mov     byte ptr ds:0D2Dh, 0 */
m->ip=0x1ce3;
su_skip_first_enter=0;
dl_write8(m,m->ds,3373,0);
goto L_11CE8;
L_11CE8: ; /* IDA 0x11CE8; dec     byte ptr ds:98A5h */
m->ip=0x1ce8;
su_skip_first_enter=0;
dl_write8(m,m->ds,39077,st_dec(m,dl_read8(m,m->ds,39077),8));
goto L_11CEC;
L_11CEC: ; /* IDA 0x11CEC; retn */
m->ip=0x1cec;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_11CED: ; /* IDA 0x11CED; cmp     byte ptr ds:0D2Dh, 0A0h */
m->ip=0x1ced;
su_skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,3373),160,8);
goto L_11CF2;
L_11CF2: ; /* IDA 0x11CF2; jb      short locret_11D1F */
m->ip=0x1cf2;
su_skip_first_enter=0;
if(m->flags&1) {goto L_11D1F;}
goto L_11CF4;
L_11CF4: ; /* IDA 0x11CF4; mov     ax, 2 */
m->ip=0x1cf4;
su_skip_first_enter=0;
m->ax=(uint16_t)(2);
goto L_11CF7;
L_11CF7: ; /* IDA 0x11CF7; call    sub_20000 */
m->ip=0x1cf7;
su_skip_first_enter=0;
mx_far(m,0x1cf7,0x1cfc,h);su_invoke_old=0;goto SU_DISPATCH;
L_11CFC: ; /* IDA 0x11CFC; mov     byte ptr ds:0D2Ch, 0 */
m->ip=0x1cfc;
su_skip_first_enter=0;
dl_write8(m,m->ds,3372,0);
goto L_11D01;
L_11D01: ; /* IDA 0x11D01; mov     byte ptr ds:0D2Dh, 0 */
m->ip=0x1d01;
su_skip_first_enter=0;
dl_write8(m,m->ds,3373,0);
goto L_11D06;
L_11D06: ; /* IDA 0x11D06; mov     byte ptr ds:98A5h, 0 */
m->ip=0x1d06;
su_skip_first_enter=0;
dl_write8(m,m->ds,39077,0);
goto L_11D0B;
L_11D0B: ; /* IDA 0x11D0B; cmp     cs:byte_10D2A, 1 */
m->ip=0x1d0b;
su_skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->cs,3370),1,8);
goto L_11D11;
L_11D11: ; /* IDA 0x11D11; jz      short loc_11D20 */
m->ip=0x1d11;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_11D20;}
goto L_11D13;
L_11D13: ; /* IDA 0x11D13; call    sub_13EFD */
m->ip=0x1d13;
su_skip_first_enter=0;
{uint16_t target=16125;ec_push(m,0x1d16);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11D16: ; /* IDA 0x11D16; call    sub_125A3 */
m->ip=0x1d16;
su_skip_first_enter=0;
{uint16_t target=9635;ec_push(m,0x1d19);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11D19: ; /* IDA 0x11D19; call    sub_12459 */
m->ip=0x1d19;
su_skip_first_enter=0;
{uint16_t target=9305;ec_push(m,0x1d1c);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11D1C: ; /* IDA 0x11D1C; call    sub_11D8E */
m->ip=0x1d1c;
su_skip_first_enter=0;
{uint16_t target=7566;ec_push(m,0x1d1f);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11D1F: ; /* IDA 0x11D1F; retn */
m->ip=0x1d1f;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_11D20: ; /* IDA 0x11D20; call    sub_10CDE */
m->ip=0x1d20;
su_skip_first_enter=0;
{uint16_t target=3294;ec_push(m,0x1d23);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11D23: ; /* IDA 0x11D23; mov     cx, 4Bh ; 'K' */
m->ip=0x1d23;
su_skip_first_enter=0;
m->cx=(uint16_t)(75);
goto L_11D26;
L_11D26: ; /* IDA 0x11D26; mov     al, 93h */
m->ip=0x1d26;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(147));
goto L_11D28;
L_11D28: ; /* IDA 0x11D28; call    sub_18810 */
m->ip=0x1d28;
su_skip_first_enter=0;
{uint16_t target=34832;ec_push(m,0x1d2b);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11D2B: ; /* IDA 0x11D2B; mov     ds, cs:word_10D52 */
m->ip=0x1d2b;
su_skip_first_enter=0;
m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
goto L_11D30;
L_11D30: ; /* IDA 0x11D30; call    sub_187FF */
m->ip=0x1d30;
su_skip_first_enter=0;
{uint16_t target=34815;ec_push(m,0x1d33);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11D33: ; /* IDA 0x11D33; mov     ah, [bx+425Eh] */
m->ip=0x1d33;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+16990))));
goto L_11D37;
L_11D37: ; /* IDA 0x11D37; mov     al, [bx+4241h] */
m->ip=0x1d37;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+16961))));
goto L_11D3B;
L_11D3B: ; /* IDA 0x11D3B; mov     cx, 197h */
m->ip=0x1d3b;
su_skip_first_enter=0;
m->cx=(uint16_t)(407);
goto L_11D3E;
L_11D3E: ; /* IDA 0x11D3E; call    sub_18810 */
m->ip=0x1d3e;
su_skip_first_enter=0;
{uint16_t target=34832;ec_push(m,0x1d41);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_11D41: ; /* IDA 0x11D41; mov     al, 2 */
m->ip=0x1d41;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(2));
goto L_11D43;
L_11D43: ; /* IDA 0x11D43; call    sub_11CB1 */
m->ip=0x1d43;
su_skip_first_enter=0;
{uint16_t target=7345;ec_push(m,0x1d46);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18B12: ; /* IDA 0x18B12; push    ax */
m->ip=0x8b12;
/* Original entry: sub_18B12 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x8b12,h->user);
ec_push(m,m->ax);
goto L_18B13;
L_18B13: ; /* IDA 0x18B13; push    bx */
m->ip=0x8b13;
su_skip_first_enter=0;
ec_push(m,m->bx);
goto L_18B14;
L_18B14: ; /* IDA 0x18B14; push    cx */
m->ip=0x8b14;
su_skip_first_enter=0;
ec_push(m,m->cx);
goto L_18B15;
L_18B15: ; /* IDA 0x18B15; push    dx */
m->ip=0x8b15;
su_skip_first_enter=0;
ec_push(m,m->dx);
goto L_18B16;
L_18B16: ; /* IDA 0x18B16; push    si */
m->ip=0x8b16;
su_skip_first_enter=0;
ec_push(m,m->si);
goto L_18B17;
L_18B17: ; /* IDA 0x18B17; push    di */
m->ip=0x8b17;
su_skip_first_enter=0;
ec_push(m,m->di);
goto L_18B18;
L_18B18: ; /* IDA 0x18B18; mov     dx, 0D9Ah */
m->ip=0x8b18;
su_skip_first_enter=0;
m->dx=(uint16_t)(3482);
goto L_18B1B;
L_18B1B: ; /* IDA 0x18B1B; xor     al, al */
m->ip=0x8b1b;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_18B1D;
L_18B1D: ; /* IDA 0x18B1D; call    sub_18B7C */
m->ip=0x8b1d;
su_skip_first_enter=0;
{uint16_t target=35708;ec_push(m,0x8b20);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18B20: ; /* IDA 0x18B20; jb      short loc_18B39 */
m->ip=0x8b20;
su_skip_first_enter=0;
if(m->flags&1) {goto L_18B39;}
goto L_18B22;
L_18B22: ; /* IDA 0x18B22; call    sub_18CAE */
m->ip=0x8b22;
su_skip_first_enter=0;
{uint16_t target=36014;ec_push(m,0x8b25);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18B25: ; /* IDA 0x18B25; push    ds */
m->ip=0x8b25;
su_skip_first_enter=0;
ec_push(m,m->ds);
goto L_18B26;
L_18B26: ; /* IDA 0x18B26; mov     ds, cs:word_10D52 */
m->ip=0x8b26;
su_skip_first_enter=0;
m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
goto L_18B2B;
L_18B2B: ; /* IDA 0x18B2B; mov     bx, cs:word_10CFD */
m->ip=0x8b2b;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->cs,3325));
goto L_18B30;
L_18B30: ; /* IDA 0x18B30; mov     al, [bx+2Bh] */
m->ip=0x8b30;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+43))));
goto L_18B33;
L_18B33: ; /* IDA 0x18B33; mov     cs:byte_10D00, al */
m->ip=0x8b33;
su_skip_first_enter=0;
dl_write8(m,m->cs,3328,(uint8_t)m->ax);
goto L_18B37;
L_18B37: ; /* IDA 0x18B37; pop     ds */
m->ip=0x8b37;
su_skip_first_enter=0;
m->ds=(uint16_t)(ec_pop(m));
goto L_18B38;
L_18B38: ; /* IDA 0x18B38; clc */
m->ip=0x8b38;
su_skip_first_enter=0;
m->flags&=(uint16_t)~1u;
goto L_18B39;
L_18B39: ; /* IDA 0x18B39; pop     di */
m->ip=0x8b39;
su_skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_18B3A;
L_18B3A: ; /* IDA 0x18B3A; pop     si */
m->ip=0x8b3a;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_pop(m));
goto L_18B3B;
L_18B3B: ; /* IDA 0x18B3B; pop     dx */
m->ip=0x8b3b;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_pop(m));
goto L_18B3C;
L_18B3C: ; /* IDA 0x18B3C; pop     cx */
m->ip=0x8b3c;
su_skip_first_enter=0;
m->cx=(uint16_t)(ec_pop(m));
goto L_18B3D;
L_18B3D: ; /* IDA 0x18B3D; pop     bx */
m->ip=0x8b3d;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_18B3E;
L_18B3E: ; /* IDA 0x18B3E; pop     ax */
m->ip=0x8b3e;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_18B3F;
L_18B3F: ; /* IDA 0x18B3F; retn */
m->ip=0x8b3f;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_18B40: ; /* IDA 0x18B40; push    ax */
m->ip=0x8b40;
/* Original entry: sub_18B40 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x8b40,h->user);
ec_push(m,m->ax);
goto L_18B41;
L_18B41: ; /* IDA 0x18B41; push    bx */
m->ip=0x8b41;
su_skip_first_enter=0;
ec_push(m,m->bx);
goto L_18B42;
L_18B42: ; /* IDA 0x18B42; push    cx */
m->ip=0x8b42;
su_skip_first_enter=0;
ec_push(m,m->cx);
goto L_18B43;
L_18B43: ; /* IDA 0x18B43; push    dx */
m->ip=0x8b43;
su_skip_first_enter=0;
ec_push(m,m->dx);
goto L_18B44;
L_18B44: ; /* IDA 0x18B44; push    si */
m->ip=0x8b44;
su_skip_first_enter=0;
ec_push(m,m->si);
goto L_18B45;
L_18B45: ; /* IDA 0x18B45; push    di */
m->ip=0x8b45;
su_skip_first_enter=0;
ec_push(m,m->di);
goto L_18B46;
L_18B46: ; /* IDA 0x18B46; mov     dx, 0DA6h */
m->ip=0x8b46;
su_skip_first_enter=0;
m->dx=(uint16_t)(3494);
goto L_18B49;
L_18B49: ; /* IDA 0x18B49; mov     al, 1 */
m->ip=0x8b49;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(1));
goto L_18B4B;
L_18B4B: ; /* IDA 0x18B4B; call    sub_18B7C */
m->ip=0x8b4b;
su_skip_first_enter=0;
{uint16_t target=35708;ec_push(m,0x8b4e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18B4E: ; /* IDA 0x18B4E; jb      short loc_18B56 */
m->ip=0x8b4e;
su_skip_first_enter=0;
if(m->flags&1) {goto L_18B56;}
goto L_18B50;
L_18B50: ; /* IDA 0x18B50; mov     ah, 1 */
m->ip=0x8b50;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(1));
goto L_18B52;
L_18B52: ; /* IDA 0x18B52; call    sub_18CAE */
m->ip=0x8b52;
su_skip_first_enter=0;
{uint16_t target=36014;ec_push(m,0x8b55);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18B55: ; /* IDA 0x18B55; clc */
m->ip=0x8b55;
su_skip_first_enter=0;
m->flags&=(uint16_t)~1u;
goto L_18B56;
L_18B56: ; /* IDA 0x18B56; pop     di */
m->ip=0x8b56;
su_skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_18B57;
L_18B57: ; /* IDA 0x18B57; pop     si */
m->ip=0x8b57;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_pop(m));
goto L_18B58;
L_18B58: ; /* IDA 0x18B58; pop     dx */
m->ip=0x8b58;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_pop(m));
goto L_18B59;
L_18B59: ; /* IDA 0x18B59; pop     cx */
m->ip=0x8b59;
su_skip_first_enter=0;
m->cx=(uint16_t)(ec_pop(m));
goto L_18B5A;
L_18B5A: ; /* IDA 0x18B5A; pop     bx */
m->ip=0x8b5a;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_18B5B;
L_18B5B: ; /* IDA 0x18B5B; pop     ax */
m->ip=0x8b5b;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_18B5C;
L_18B5C: ; /* IDA 0x18B5C; retn */
m->ip=0x8b5c;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_18CAE: ; /* IDA 0x18CAE; push    ax */
m->ip=0x8cae;
/* Original entry: sub_18CAE */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x8cae,h->user);
ec_push(m,m->ax);
goto L_18CAF;
L_18CAF: ; /* IDA 0x18CAF; mov     bx, 0D9Ah */
m->ip=0x8caf;
su_skip_first_enter=0;
m->bx=(uint16_t)(3482);
goto L_18CB2;
L_18CB2: ; /* IDA 0x18CB2; and     ah, ah */
m->ip=0x8cb2;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->ax>>8))&((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
goto L_18CB4;
L_18CB4: ; /* IDA 0x18CB4; jz      short loc_18CB9 */
m->ip=0x8cb4;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_18CB9;}
goto L_18CB6;
L_18CB6: ; /* IDA 0x18CB6; mov     bx, 0DA6h */
m->ip=0x8cb6;
su_skip_first_enter=0;
m->bx=(uint16_t)(3494);
goto L_18CB9;
L_18CB9: ; /* IDA 0x18CB9; xor     ah, ah */
m->ip=0x8cb9;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
goto L_18CBB;
L_18CBB: ; /* IDA 0x18CBB; mov     dx, 56C0h */
m->ip=0x8cbb;
su_skip_first_enter=0;
m->dx=(uint16_t)(22208);
goto L_18CBE;
L_18CBE: ; /* IDA 0x18CBE; mul     dx */
m->ip=0x8cbe;
su_skip_first_enter=0;
dl_mul(m,m->dx,16);
goto L_18CC0;
L_18CC0: ; /* IDA 0x18CC0; mov     cx, dx */
m->ip=0x8cc0;
su_skip_first_enter=0;
m->cx=(uint16_t)(m->dx);
goto L_18CC2;
L_18CC2: ; /* IDA 0x18CC2; mov     dx, bx */
m->ip=0x8cc2;
su_skip_first_enter=0;
m->dx=(uint16_t)(m->bx);
goto L_18CC4;
L_18CC4: ; /* IDA 0x18CC4; mov     bx, cs */
m->ip=0x8cc4;
su_skip_first_enter=0;
m->bx=(uint16_t)(m->cs);
goto L_18CC6;
L_18CC6: ; /* IDA 0x18CC6; mov     si, 0CF0h */
m->ip=0x8cc6;
su_skip_first_enter=0;
m->si=(uint16_t)(3312);
goto L_18CC9;
L_18CC9: ; /* IDA 0x18CC9; mov     di, 3Bh ; ';' */
m->ip=0x8cc9;
su_skip_first_enter=0;
m->di=(uint16_t)(59);
goto L_18CCC;
L_18CCC: ; /* IDA 0x18CCC; call    sub_1E38C */
m->ip=0x8ccc;
su_skip_first_enter=0;
{uint16_t target=58252;ec_push(m,0x8ccf);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18CCF: ; /* IDA 0x18CCF; add     ax, 80h */
m->ip=0x8ccf;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_math(m,m->ax,128,16,0,0));
goto L_18CD2;
L_18CD2: ; /* IDA 0x18CD2; adc     cx, 0 */
m->ip=0x8cd2;
su_skip_first_enter=0;
m->cx=(uint16_t)(ec_math(m,m->cx,0,16,m->flags&1,0));
goto L_18CD5;
L_18CD5: ; /* IDA 0x18CD5; mov     bx, ds:0D52h */
m->ip=0x8cd5;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,3410));
goto L_18CD9;
L_18CD9: ; /* IDA 0x18CD9; xor     si, si */
m->ip=0x8cd9;
su_skip_first_enter=0;
{uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_18CDB;
L_18CDB: ; /* IDA 0x18CDB; mov     di, 5240h */
m->ip=0x8cdb;
su_skip_first_enter=0;
m->di=(uint16_t)(21056);
goto L_18CDE;
L_18CDE: ; /* IDA 0x18CDE; call    sub_1E38C */
m->ip=0x8cde;
su_skip_first_enter=0;
{uint16_t target=58252;ec_push(m,0x8ce1);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18CE1: ; /* IDA 0x18CE1; add     ax, 5240h */
m->ip=0x8ce1;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_math(m,m->ax,21056,16,0,0));
goto L_18CE4;
L_18CE4: ; /* IDA 0x18CE4; adc     cx, 0 */
m->ip=0x8ce4;
su_skip_first_enter=0;
m->cx=(uint16_t)(ec_math(m,m->cx,0,16,m->flags&1,0));
goto L_18CE7;
L_18CE7: ; /* IDA 0x18CE7; mov     bx, ds:0D56h */
m->ip=0x8ce7;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_read16(m,m->ds,3414));
goto L_18CEB;
L_18CEB: ; /* IDA 0x18CEB; xor     si, si */
m->ip=0x8ceb;
su_skip_first_enter=0;
{uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_18CED;
L_18CED: ; /* IDA 0x18CED; mov     di, 400h */
m->ip=0x8ced;
su_skip_first_enter=0;
m->di=(uint16_t)(1024);
goto L_18CF0;
L_18CF0: ; /* IDA 0x18CF0; call    sub_1E38C */
m->ip=0x8cf0;
su_skip_first_enter=0;
{uint16_t target=58252;ec_push(m,0x8cf3);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18CF3: ; /* IDA 0x18CF3; pop     ax */
m->ip=0x8cf3;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_18CF4;
L_18CF4: ; /* IDA 0x18CF4; and     ah, ah */
m->ip=0x8cf4;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->ax>>8))&((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
goto L_18CF6;
L_18CF6: ; /* IDA 0x18CF6; jnz     short locret_18CFE */
m->ip=0x8cf6;
su_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_18CFE;}
goto L_18CF8;
L_18CF8: ; /* IDA 0x18CF8; mov     cs:byte_10CF8, 0 */
m->ip=0x8cf8;
su_skip_first_enter=0;
dl_write8(m,m->cs,3320,0);
goto L_18CFE;
L_18CFE: ; /* IDA 0x18CFE; retn */
m->ip=0x8cfe;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_18DC8: ; /* IDA 0x18DC8; push    ds */
m->ip=0x8dc8;
/* Original entry: sub_18DC8 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x8dc8,h->user);
ec_push(m,m->ds);
goto L_18DC9;
L_18DC9: ; /* IDA 0x18DC9; push    cx */
m->ip=0x8dc9;
su_skip_first_enter=0;
ec_push(m,m->cx);
goto L_18DCA;
L_18DCA: ; /* IDA 0x18DCA; push    bx */
m->ip=0x8dca;
su_skip_first_enter=0;
ec_push(m,m->bx);
goto L_18DCB;
L_18DCB: ; /* IDA 0x18DCB; push    dx */
m->ip=0x8dcb;
su_skip_first_enter=0;
ec_push(m,m->dx);
goto L_18DCC;
L_18DCC: ; /* IDA 0x18DCC; push    si */
m->ip=0x8dcc;
su_skip_first_enter=0;
ec_push(m,m->si);
goto L_18DCD;
L_18DCD: ; /* IDA 0x18DCD; push    bp */
m->ip=0x8dcd;
su_skip_first_enter=0;
ec_push(m,m->bp);
goto L_18DCE;
L_18DCE: ; /* IDA 0x18DCE; sub     sp, 4 */
m->ip=0x8dce;
su_skip_first_enter=0;
m->sp=(uint16_t)(ec_math(m,m->sp,4,16,0,1));
goto L_18DD1;
L_18DD1: ; /* IDA 0x18DD1; mov     bp, sp */
m->ip=0x8dd1;
su_skip_first_enter=0;
m->bp=(uint16_t)(m->sp);
goto L_18DD3;
L_18DD3: ; /* IDA 0x18DD3; call    sub_101B4 */
m->ip=0x8dd3;
su_skip_first_enter=0;
{uint16_t target=436;ec_push(m,0x8dd6);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18DD6: ; /* IDA 0x18DD6; mov     [bp+0], dx */
m->ip=0x8dd6;
su_skip_first_enter=0;
dl_write16(m,m->ss,(uint16_t)(m->bp),m->dx);
goto L_18DD9;
L_18DD9: ; /* IDA 0x18DD9; mov     [bp+2], bx */
m->ip=0x8dd9;
su_skip_first_enter=0;
dl_write16(m,m->ss,(uint16_t)(m->bp+2),m->bx);
goto L_18DDC;
L_18DDC: ; /* IDA 0x18DDC; mov     cx, 600Dh */
m->ip=0x8ddc;
su_skip_first_enter=0;
m->cx=(uint16_t)(24589);
goto L_18DDF;
L_18DDF: ; /* IDA 0x18DDF; call    sub_19796 */
m->ip=0x8ddf;
su_skip_first_enter=0;
{uint16_t target=38806;ec_push(m,0x8de2);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18DE2: ; /* IDA 0x18DE2; mov     cl, 4 */
m->ip=0x8de2;
su_skip_first_enter=0;
vg_cl(m,(uint8_t)(4));
goto L_18DE4;
L_18DE4: ; /* IDA 0x18DE4; push    bx */
m->ip=0x8de4;
su_skip_first_enter=0;
ec_push(m,m->bx);
goto L_18DE5;
L_18DE5: ; /* IDA 0x18DE5; push    dx */
m->ip=0x8de5;
su_skip_first_enter=0;
ec_push(m,m->dx);
goto L_18DE6;
L_18DE6: ; /* IDA 0x18DE6; shr     dx, cl */
m->ip=0x8de6;
su_skip_first_enter=0;
m->dx=(uint16_t)(dl_shift(m,m->dx,(uint8_t)m->cx,16,5));
goto L_18DE8;
L_18DE8: ; /* IDA 0x18DE8; shr     bx, cl */
m->ip=0x8de8;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,(uint8_t)m->cx,16,5));
goto L_18DEA;
L_18DEA: ; /* IDA 0x18DEA; mov     cx, 60Dh */
m->ip=0x8dea;
su_skip_first_enter=0;
m->cx=(uint16_t)(1549);
goto L_18DED;
L_18DED: ; /* IDA 0x18DED; call    sub_10C14 */
m->ip=0x8ded;
su_skip_first_enter=0;
{uint16_t target=3092;ec_push(m,0x8df0);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18DF0: ; /* IDA 0x18DF0; pop     dx */
m->ip=0x8df0;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_pop(m));
goto L_18DF1;
L_18DF1: ; /* IDA 0x18DF1; pop     bx */
m->ip=0x8df1;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_18DF2;
L_18DF2: ; /* IDA 0x18DF2; add     dx, 8 */
m->ip=0x8df2;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_math(m,m->dx,8,16,0,0));
goto L_18DF5;
L_18DF5: ; /* IDA 0x18DF5; add     bx, 8 */
m->ip=0x8df5;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,8,16,0,0));
goto L_18DF8;
L_18DF8: ; /* IDA 0x18DF8; mov     al, 7 */
m->ip=0x8df8;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(7));
goto L_18DFA;
L_18DFA: ; /* IDA 0x18DFA; call    sub_10337 */
m->ip=0x8dfa;
su_skip_first_enter=0;
{uint16_t target=823;ec_push(m,0x8dfd);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18DFD: ; /* IDA 0x18DFD; sub     dx, 8 */
m->ip=0x8dfd;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_math(m,m->dx,8,16,0,1));
goto L_18E00;
L_18E00: ; /* IDA 0x18E00; sub     bx, 8 */
m->ip=0x8e00;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,8,16,0,1));
goto L_18E03;
L_18E03: ; /* IDA 0x18E03; add     bx, 0Ch */
m->ip=0x8e03;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,12,16,0,0));
goto L_18E06;
L_18E06: ; /* IDA 0x18E06; add     dx, di */
m->ip=0x8e06;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_math(m,m->dx,m->di,16,0,0));
goto L_18E08;
L_18E08: ; /* IDA 0x18E08; mov     ax, 0F03h */
m->ip=0x8e08;
su_skip_first_enter=0;
m->ax=(uint16_t)(3843);
goto L_18E0B;
L_18E0B: ; /* IDA 0x18E0B; call    sub_106F5 */
m->ip=0x8e0b;
su_skip_first_enter=0;
{uint16_t target=1781;ec_push(m,0x8e0e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E0E: ; /* IDA 0x18E0E; call    sub_101DB */
m->ip=0x8e0e;
su_skip_first_enter=0;
{uint16_t target=475;ec_push(m,0x8e11);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E11: ; /* IDA 0x18E11; add     word ptr [bp+0], 28h ; '(' */
m->ip=0x8e11;
su_skip_first_enter=0;
dl_write16(m,m->ss,(uint16_t)(m->bp),ec_math(m,dl_read16(m,m->ss,(uint16_t)(m->bp)),40,16,0,0));
goto L_18E15;
L_18E15: ; /* IDA 0x18E15; add     word ptr [bp+2], 28h ; '(' */
m->ip=0x8e15;
su_skip_first_enter=0;
dl_write16(m,m->ss,(uint16_t)(m->bp+2),ec_math(m,dl_read16(m,m->ss,(uint16_t)(m->bp+2)),40,16,0,0));
goto L_18E19;
L_18E19: ; /* IDA 0x18E19; call    sub_121E7 */
m->ip=0x8e19;
su_skip_first_enter=0;
{uint16_t target=8679;ec_push(m,0x8e1c);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E1C: ; /* IDA 0x18E1C; jb      short loc_18E19 */
m->ip=0x8e1c;
su_skip_first_enter=0;
if(m->flags&1) {goto L_18E19;}
goto L_18E1E;
L_18E1E: ; /* IDA 0x18E1E; sub     cx, [bp+0] */
m->ip=0x8e1e;
su_skip_first_enter=0;
m->cx=(uint16_t)(ec_math(m,m->cx,dl_read16(m,m->ss,(uint16_t)(m->bp)),16,0,1));
goto L_18E21;
L_18E21: ; /* IDA 0x18E21; jb      short loc_18E19 */
m->ip=0x8e21;
su_skip_first_enter=0;
if(m->flags&1) {goto L_18E19;}
goto L_18E23;
L_18E23: ; /* IDA 0x18E23; cmp     cx, 80h */
m->ip=0x8e23;
su_skip_first_enter=0;
ec_cmp(m,m->cx,128,16);
goto L_18E27;
L_18E27: ; /* IDA 0x18E27; jnb     short loc_18E19 */
m->ip=0x8e27;
su_skip_first_enter=0;
if(!(m->flags&1)) {goto L_18E19;}
goto L_18E29;
L_18E29: ; /* IDA 0x18E29; sub     dx, [bp+2] */
m->ip=0x8e29;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_math(m,m->dx,dl_read16(m,m->ss,(uint16_t)(m->bp+2)),16,0,1));
goto L_18E2C;
L_18E2C: ; /* IDA 0x18E2C; jb      short loc_18E19 */
m->ip=0x8e2c;
su_skip_first_enter=0;
if(m->flags&1) {goto L_18E19;}
goto L_18E2E;
L_18E2E: ; /* IDA 0x18E2E; cmp     dx, 48h ; 'H' */
m->ip=0x8e2e;
su_skip_first_enter=0;
ec_cmp(m,m->dx,72,16);
goto L_18E31;
L_18E31: ; /* IDA 0x18E31; jnb     short loc_18E19 */
m->ip=0x8e31;
su_skip_first_enter=0;
if(!(m->flags&1)) {goto L_18E19;}
goto L_18E33;
L_18E33: ; /* IDA 0x18E33; mov     ax, 0 */
m->ip=0x8e33;
su_skip_first_enter=0;
m->ax=(uint16_t)(0);
goto L_18E36;
L_18E36: ; /* IDA 0x18E36; shr     dx, 1 */
m->ip=0x8e36;
su_skip_first_enter=0;
m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
goto L_18E38;
L_18E38: ; /* IDA 0x18E38; shr     dx, 1 */
m->ip=0x8e38;
su_skip_first_enter=0;
m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
goto L_18E3A;
L_18E3A: ; /* IDA 0x18E3A; shr     dx, 1 */
m->ip=0x8e3a;
su_skip_first_enter=0;
m->dx=(uint16_t)(dl_shift(m,m->dx,1,16,5));
goto L_18E3C;
L_18E3C: ; /* IDA 0x18E3C; cmp     dx, 2 */
m->ip=0x8e3c;
su_skip_first_enter=0;
ec_cmp(m,m->dx,2,16);
goto L_18E3F;
L_18E3F: ; /* IDA 0x18E3F; jz      short loc_18E19 */
m->ip=0x8e3f;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_18E19;}
goto L_18E41;
L_18E41: ; /* IDA 0x18E41; jb      short loc_18E44 */
m->ip=0x8e41;
su_skip_first_enter=0;
if(m->flags&1) {goto L_18E44;}
goto L_18E43;
L_18E43: ; /* IDA 0x18E43; inc     ax */
m->ip=0x8e43;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_inc(m,m->ax,16));
goto L_18E44;
L_18E44: ; /* IDA 0x18E44; add     sp, 4 */
m->ip=0x8e44;
su_skip_first_enter=0;
m->sp=(uint16_t)(ec_math(m,m->sp,4,16,0,0));
goto L_18E47;
L_18E47: ; /* IDA 0x18E47; pop     bp */
m->ip=0x8e47;
su_skip_first_enter=0;
m->bp=(uint16_t)(ec_pop(m));
goto L_18E48;
L_18E48: ; /* IDA 0x18E48; pop     si */
m->ip=0x8e48;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_pop(m));
goto L_18E49;
L_18E49: ; /* IDA 0x18E49; pop     dx */
m->ip=0x8e49;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_pop(m));
goto L_18E4A;
L_18E4A: ; /* IDA 0x18E4A; pop     bx */
m->ip=0x8e4a;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_18E4B;
L_18E4B: ; /* IDA 0x18E4B; call    sub_101B4 */
m->ip=0x8e4b;
su_skip_first_enter=0;
{uint16_t target=436;ec_push(m,0x8e4e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E4E: ; /* IDA 0x18E4E; mov     cx, 600Dh */
m->ip=0x8e4e;
su_skip_first_enter=0;
m->cx=(uint16_t)(24589);
goto L_18E51;
L_18E51: ; /* IDA 0x18E51; call    sub_197C3 */
m->ip=0x8e51;
su_skip_first_enter=0;
{uint16_t target=38851;ec_push(m,0x8e54);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E54: ; /* IDA 0x18E54; pop     cx */
m->ip=0x8e54;
su_skip_first_enter=0;
m->cx=(uint16_t)(ec_pop(m));
goto L_18E55;
L_18E55: ; /* IDA 0x18E55; pop     ds */
m->ip=0x8e55;
su_skip_first_enter=0;
m->ds=(uint16_t)(ec_pop(m));
goto L_18E56;
L_18E56: ; /* IDA 0x18E56; call    sub_101DB */
m->ip=0x8e56;
su_skip_first_enter=0;
{uint16_t target=475;ec_push(m,0x8e59);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E59: ; /* IDA 0x18E59; retn */
m->ip=0x8e59;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_18E5A: ; /* IDA 0x18E5A; push    ds */
m->ip=0x8e5a;
/* Original entry: sub_18E5A */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x8e5a,h->user);
ec_push(m,m->ds);
goto L_18E5B;
L_18E5B: ; /* IDA 0x18E5B; push    bx */
m->ip=0x8e5b;
su_skip_first_enter=0;
ec_push(m,m->bx);
goto L_18E5C;
L_18E5C: ; /* IDA 0x18E5C; push    si */
m->ip=0x8e5c;
su_skip_first_enter=0;
ec_push(m,m->si);
goto L_18E5D;
L_18E5D: ; /* IDA 0x18E5D; mov     ds, cs:word_10D52 */
m->ip=0x8e5d;
su_skip_first_enter=0;
m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
goto L_18E62;
L_18E62: ; /* IDA 0x18E62; mov     si, bx */
m->ip=0x8e62;
su_skip_first_enter=0;
m->si=(uint16_t)(m->bx);
goto L_18E64;
L_18E64: ; /* IDA 0x18E64; call    sub_101B4 */
m->ip=0x8e64;
su_skip_first_enter=0;
{uint16_t target=436;ec_push(m,0x8e67);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E67: ; /* IDA 0x18E67; call    sub_18F37 */
m->ip=0x8e67;
su_skip_first_enter=0;
{uint16_t target=36663;ec_push(m,0x8e6a);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E6A: ; /* IDA 0x18E6A; call    sub_18EA0 */
m->ip=0x8e6a;
su_skip_first_enter=0;
{uint16_t target=36512;ec_push(m,0x8e6d);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E6D: ; /* IDA 0x18E6D; call    sub_101DB */
m->ip=0x8e6d;
su_skip_first_enter=0;
{uint16_t target=475;ec_push(m,0x8e70);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E70: ; /* IDA 0x18E70; call    sub_121E7 */
m->ip=0x8e70;
su_skip_first_enter=0;
{uint16_t target=8679;ec_push(m,0x8e73);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E73: ; /* IDA 0x18E73; jb      short loc_18E91 */
m->ip=0x8e73;
su_skip_first_enter=0;
if(m->flags&1) {goto L_18E91;}
goto L_18E75;
L_18E75: ; /* IDA 0x18E75; call    sub_1E453 */
m->ip=0x8e75;
su_skip_first_enter=0;
{uint16_t target=58451;ec_push(m,0x8e78);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E78: ; /* IDA 0x18E78; sub     al, 20h ; ' ' */
m->ip=0x8e78;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,32,8,0,1)));
goto L_18E7A;
L_18E7A: ; /* IDA 0x18E7A; jb      short loc_18E70 */
m->ip=0x8e7a;
su_skip_first_enter=0;
if(m->flags&1) {goto L_18E70;}
goto L_18E7C;
L_18E7C: ; /* IDA 0x18E7C; ja      short loc_18E8E */
m->ip=0x8e7c;
su_skip_first_enter=0;
if(!(m->flags&(1|0x40))) {goto L_18E8E;}
goto L_18E7E;
L_18E7E: ; /* IDA 0x18E7E; cmp     byte ptr [si+2], 7Fh */
m->ip=0x8e7e;
su_skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,(uint16_t)(m->si+2)),127,8);
goto L_18E82;
L_18E82: ; /* IDA 0x18E82; jnz     short loc_18E8E */
m->ip=0x8e82;
su_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_18E8E;}
goto L_18E84;
L_18E84: ; /* IDA 0x18E84; cmp     word ptr ds:5222h, 0D0A1h */
m->ip=0x8e84;
su_skip_first_enter=0;
ec_cmp(m,dl_read16(m,m->ds,21026),53409,16);
goto L_18E8A;
L_18E8A: ; /* IDA 0x18E8A; jnz     short loc_18E8E */
m->ip=0x8e8a;
su_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_18E8E;}
goto L_18E8C;
L_18E8C: ; /* IDA 0x18E8C; mov     al, 1 */
m->ip=0x8e8c;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(1));
goto L_18E8E;
L_18E8E: ; /* IDA 0x18E8E; clc */
m->ip=0x8e8e;
su_skip_first_enter=0;
m->flags&=(uint16_t)~1u;
goto L_18E8F;
L_18E8F: ; /* IDA 0x18E8F; jmp     short loc_18E95 */
m->ip=0x8e8f;
su_skip_first_enter=0;
goto L_18E95;
L_18E91: ; /* IDA 0x18E91; call    sub_18F7C */
m->ip=0x8e91;
su_skip_first_enter=0;
{uint16_t target=36732;ec_push(m,0x8e94);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E94: ; /* IDA 0x18E94; stc */
m->ip=0x8e94;
su_skip_first_enter=0;
m->flags|=1;
goto L_18E95;
L_18E95: ; /* IDA 0x18E95; pushf */
m->ip=0x8e95;
su_skip_first_enter=0;
ec_push(m,m->flags);
goto L_18E96;
L_18E96: ; /* IDA 0x18E96; push    ax */
m->ip=0x8e96;
su_skip_first_enter=0;
ec_push(m,m->ax);
goto L_18E97;
L_18E97: ; /* IDA 0x18E97; call    sub_18F6D */
m->ip=0x8e97;
su_skip_first_enter=0;
{uint16_t target=36717;ec_push(m,0x8e9a);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18E9A: ; /* IDA 0x18E9A; pop     ax */
m->ip=0x8e9a;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_18E9B;
L_18E9B: ; /* IDA 0x18E9B; popf */
m->ip=0x8e9b;
su_skip_first_enter=0;
m->flags=ec_pop(m)|2;
goto L_18E9C;
L_18E9C: ; /* IDA 0x18E9C; pop     si */
m->ip=0x8e9c;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_pop(m));
goto L_18E9D;
L_18E9D: ; /* IDA 0x18E9D; pop     bx */
m->ip=0x8e9d;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_18E9E;
L_18E9E: ; /* IDA 0x18E9E; pop     ds */
m->ip=0x8e9e;
su_skip_first_enter=0;
m->ds=(uint16_t)(ec_pop(m));
goto L_18E9F;
L_18E9F: ; /* IDA 0x18E9F; retn */
m->ip=0x8e9f;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_18EA0: ; /* IDA 0x18EA0; push    si */
m->ip=0x8ea0;
/* Original entry: sub_18EA0 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x8ea0,h->user);
ec_push(m,m->si);
goto L_18EA1;
L_18EA1: ; /* IDA 0x18EA1; push    di */
m->ip=0x8ea1;
su_skip_first_enter=0;
ec_push(m,m->di);
goto L_18EA2;
L_18EA2: ; /* IDA 0x18EA2; mov     di, si */
m->ip=0x8ea2;
su_skip_first_enter=0;
m->di=(uint16_t)(m->si);
goto L_18EA4;
L_18EA4: ; /* IDA 0x18EA4; mov     bh, [si+1] */
m->ip=0x8ea4;
su_skip_first_enter=0;
gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+1))));
goto L_18EA7;
L_18EA7: ; /* IDA 0x18EA7; xor     bl, bl */
m->ip=0x8ea7;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
goto L_18EA9;
L_18EA9: ; /* IDA 0x18EA9; shr     bx, 1 */
m->ip=0x8ea9;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_18EAB;
L_18EAB: ; /* IDA 0x18EAB; shr     bx, 1 */
m->ip=0x8eab;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_18EAD;
L_18EAD: ; /* IDA 0x18EAD; shr     bx, 1 */
m->ip=0x8ead;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_18EAF;
L_18EAF: ; /* IDA 0x18EAF; mov     si, bx */
m->ip=0x8eaf;
su_skip_first_enter=0;
m->si=(uint16_t)(m->bx);
goto L_18EB1;
L_18EB1: ; /* IDA 0x18EB1; mov     al, [bx+4241h] */
m->ip=0x8eb1;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+16961))));
goto L_18EB5;
L_18EB5: ; /* IDA 0x18EB5; mov     bx, 2817h */
m->ip=0x8eb5;
su_skip_first_enter=0;
m->bx=(uint16_t)(10263);
goto L_18EB8;
L_18EB8: ; /* IDA 0x18EB8; call    sub_107D2 */
m->ip=0x8eb8;
su_skip_first_enter=0;
{uint16_t target=2002;ec_push(m,0x8ebb);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18EBB: ; /* IDA 0x18EBB; add     si, 4242h */
m->ip=0x8ebb;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_math(m,m->si,16962,16,0,0));
goto L_18EBF;
L_18EBF: ; /* IDA 0x18EBF; mov     ax, 0F03h */
m->ip=0x8ebf;
su_skip_first_enter=0;
m->ax=(uint16_t)(3843);
goto L_18EC2;
L_18EC2: ; /* IDA 0x18EC2; mov     dx, 0C8h */
m->ip=0x8ec2;
su_skip_first_enter=0;
m->dx=(uint16_t)(200);
goto L_18EC5;
L_18EC5: ; /* IDA 0x18EC5; mov     bx, 0D8h */
m->ip=0x8ec5;
su_skip_first_enter=0;
m->bx=(uint16_t)(216);
goto L_18EC8;
L_18EC8: ; /* IDA 0x18EC8; call    sub_106FD */
m->ip=0x8ec8;
su_skip_first_enter=0;
{uint16_t target=1789;ec_push(m,0x8ecb);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18ECB: ; /* IDA 0x18ECB; mov     bh, [di+2] */
m->ip=0x8ecb;
su_skip_first_enter=0;
gw_bh(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+2))));
goto L_18ECE;
L_18ECE: ; /* IDA 0x18ECE; xor     bl, bl */
m->ip=0x8ece;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->bx)^((uint8_t)m->bx))&255;ec_logic(m,result,8);gw_bl(m,(uint8_t)(result));}
goto L_18ED0;
L_18ED0: ; /* IDA 0x18ED0; shr     bx, 1 */
m->ip=0x8ed0;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_18ED2;
L_18ED2: ; /* IDA 0x18ED2; shr     bx, 1 */
m->ip=0x8ed2;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_18ED4;
L_18ED4: ; /* IDA 0x18ED4; shr     bx, 1 */
m->ip=0x8ed4;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,5));
goto L_18ED6;
L_18ED6: ; /* IDA 0x18ED6; mov     si, bx */
m->ip=0x8ed6;
su_skip_first_enter=0;
m->si=(uint16_t)(m->bx);
goto L_18ED8;
L_18ED8: ; /* IDA 0x18ED8; mov     al, [bx+4241h] */
m->ip=0x8ed8;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->bx+16961))));
goto L_18EDC;
L_18EDC: ; /* IDA 0x18EDC; cmp     al, 0FFh */
m->ip=0x8edc;
su_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,255,8);
goto L_18EDE;
L_18EDE: ; /* IDA 0x18EDE; jz      short loc_18EE6 */
m->ip=0x8ede;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_18EE6;}
goto L_18EE0;
L_18EE0: ; /* IDA 0x18EE0; mov     bx, 34A7h */
m->ip=0x8ee0;
su_skip_first_enter=0;
m->bx=(uint16_t)(13479);
goto L_18EE3;
L_18EE3: ; /* IDA 0x18EE3; call    sub_107D2 */
m->ip=0x8ee3;
su_skip_first_enter=0;
{uint16_t target=2002;ec_push(m,0x8ee6);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18EE6: ; /* IDA 0x18EE6; add     si, 4242h */
m->ip=0x8ee6;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_math(m,m->si,16962,16,0,0));
goto L_18EEA;
L_18EEA: ; /* IDA 0x18EEA; mov     ax, 0F03h */
m->ip=0x8eea;
su_skip_first_enter=0;
m->ax=(uint16_t)(3843);
goto L_18EED;
L_18EED: ; /* IDA 0x18EED; mov     dx, 148h */
m->ip=0x8eed;
su_skip_first_enter=0;
m->dx=(uint16_t)(328);
goto L_18EF0;
L_18EF0: ; /* IDA 0x18EF0; mov     bx, 90h */
m->ip=0x8ef0;
su_skip_first_enter=0;
m->bx=(uint16_t)(144);
goto L_18EF3;
L_18EF3: ; /* IDA 0x18EF3; call    sub_106FD */
m->ip=0x8ef3;
su_skip_first_enter=0;
{uint16_t target=1789;ec_push(m,0x8ef6);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18EF6: ; /* IDA 0x18EF6; mov     ah, [di+3] */
m->ip=0x8ef6;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->di+3))));
goto L_18EF9;
L_18EF9: ; /* IDA 0x18EF9; xor     al, al */
m->ip=0x8ef9;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_18EFB;
L_18EFB: ; /* IDA 0x18EFB; shr     ax, 1 */
m->ip=0x8efb;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
goto L_18EFD;
L_18EFD: ; /* IDA 0x18EFD; shr     ax, 1 */
m->ip=0x8efd;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
goto L_18EFF;
L_18EFF: ; /* IDA 0x18EFF; shr     ax, 1 */
m->ip=0x8eff;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,5));
goto L_18F01;
L_18F01: ; /* IDA 0x18F01; add     ax, 842h */
m->ip=0x8f01;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_math(m,m->ax,2114,16,0,0));
goto L_18F04;
L_18F04: ; /* IDA 0x18F04; mov     si, ax */
m->ip=0x8f04;
su_skip_first_enter=0;
m->si=(uint16_t)(m->ax);
goto L_18F06;
L_18F06: ; /* IDA 0x18F06; mov     ax, 0F03h */
m->ip=0x8f06;
su_skip_first_enter=0;
m->ax=(uint16_t)(3843);
goto L_18F09;
L_18F09: ; /* IDA 0x18F09; mov     dx, 108h */
m->ip=0x8f09;
su_skip_first_enter=0;
m->dx=(uint16_t)(264);
goto L_18F0C;
L_18F0C: ; /* IDA 0x18F0C; mov     bx, 0F0h */
m->ip=0x8f0c;
su_skip_first_enter=0;
m->bx=(uint16_t)(240);
goto L_18F0F;
L_18F0F: ; /* IDA 0x18F0F; call    sub_106FD */
m->ip=0x8f0f;
su_skip_first_enter=0;
{uint16_t target=1789;ec_push(m,0x8f12);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18F12: ; /* IDA 0x18F12; mov     si, di */
m->ip=0x8f12;
su_skip_first_enter=0;
m->si=(uint16_t)(m->di);
goto L_18F14;
L_18F14: ; /* IDA 0x18F14; mov     al, [si+18h] */
m->ip=0x8f14;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+24))));
goto L_18F17;
L_18F17: ; /* IDA 0x18F17; xor     dx, dx */
m->ip=0x8f17;
su_skip_first_enter=0;
{uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
goto L_18F19;
L_18F19: ; /* IDA 0x18F19; mov     ah, dl */
m->ip=0x8f19;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)((uint8_t)m->dx));
goto L_18F1B;
L_18F1B: ; /* IDA 0x18F1B; mov     bx, 0F03h */
m->ip=0x8f1b;
su_skip_first_enter=0;
m->bx=(uint16_t)(3843);
goto L_18F1E;
L_18F1E: ; /* IDA 0x18F1E; mov     di, 5022h */
m->ip=0x8f1e;
su_skip_first_enter=0;
m->di=(uint16_t)(20514);
goto L_18F21;
L_18F21: ; /* IDA 0x18F21; call    sub_1062F */
m->ip=0x8f21;
su_skip_first_enter=0;
{uint16_t target=1583;ec_push(m,0x8f24);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18F24: ; /* IDA 0x18F24; mov     al, [si+23h] */
m->ip=0x8f24;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+35))));
goto L_18F27;
L_18F27: ; /* IDA 0x18F27; xor     dx, dx */
m->ip=0x8f27;
su_skip_first_enter=0;
{uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
goto L_18F29;
L_18F29: ; /* IDA 0x18F29; mov     ah, dl */
m->ip=0x8f29;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)((uint8_t)m->dx));
goto L_18F2B;
L_18F2B: ; /* IDA 0x18F2B; mov     bx, 0F03h */
m->ip=0x8f2b;
su_skip_first_enter=0;
m->bx=(uint16_t)(3843);
goto L_18F2E;
L_18F2E: ; /* IDA 0x18F2E; mov     di, 5522h */
m->ip=0x8f2e;
su_skip_first_enter=0;
m->di=(uint16_t)(21794);
goto L_18F31;
L_18F31: ; /* IDA 0x18F31; call    sub_1062F */
m->ip=0x8f31;
su_skip_first_enter=0;
{uint16_t target=1583;ec_push(m,0x8f34);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18F34: ; /* IDA 0x18F34; pop     di */
m->ip=0x8f34;
su_skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_18F35;
L_18F35: ; /* IDA 0x18F35; pop     si */
m->ip=0x8f35;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_pop(m));
goto L_18F36;
L_18F36: ; /* IDA 0x18F36; retn */
m->ip=0x8f36;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_18F37: ; /* IDA 0x18F37; mov     al, 1Fh */
m->ip=0x8f37;
/* Original entry: sub_18F37 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x8f37,h->user);
vg_al(m,(uint8_t)(31));
goto L_18F39;
L_18F39: ; /* IDA 0x18F39; mov     dx, 0Ah */
m->ip=0x8f39;
su_skip_first_enter=0;
m->dx=(uint16_t)(10);
goto L_18F3C;
L_18F3C: ; /* IDA 0x18F3C; mov     bx, 5 */
m->ip=0x8f3c;
su_skip_first_enter=0;
m->bx=(uint16_t)(5);
goto L_18F3F;
L_18F3F: ; /* IDA 0x18F3F; mov     cx, 0C0Fh */
m->ip=0x8f3f;
su_skip_first_enter=0;
m->cx=(uint16_t)(3087);
goto L_18F42;
L_18F42: ; /* IDA 0x18F42; call    sub_1895D */
m->ip=0x8f42;
su_skip_first_enter=0;
{uint16_t target=35165;ec_push(m,0x8f45);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18F45: ; /* IDA 0x18F45; mov     al, 8 */
m->ip=0x8f45;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(8));
goto L_18F47;
L_18F47: ; /* IDA 0x18F47; mov     dx, 0A8h */
m->ip=0x8f47;
su_skip_first_enter=0;
m->dx=(uint16_t)(168);
goto L_18F4A;
L_18F4A: ; /* IDA 0x18F4A; mov     bx, 78h ; 'x' */
m->ip=0x8f4a;
su_skip_first_enter=0;
m->bx=(uint16_t)(120);
goto L_18F4D;
L_18F4D: ; /* IDA 0x18F4D; call    sub_10337 */
m->ip=0x8f4d;
su_skip_first_enter=0;
{uint16_t target=823;ec_push(m,0x8f50);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18F50: ; /* IDA 0x18F50; mov     al, 20h ; ' ' */
m->ip=0x8f50;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(32));
goto L_18F52;
L_18F52: ; /* IDA 0x18F52; mov     dx, 148h */
m->ip=0x8f52;
su_skip_first_enter=0;
m->dx=(uint16_t)(328);
goto L_18F55;
L_18F55: ; /* IDA 0x18F55; mov     bx, 110h */
m->ip=0x8f55;
su_skip_first_enter=0;
m->bx=(uint16_t)(272);
goto L_18F58;
L_18F58: ; /* IDA 0x18F58; mov     cx, 206h */
m->ip=0x8f58;
su_skip_first_enter=0;
m->cx=(uint16_t)(518);
goto L_18F5B;
L_18F5B: ; /* IDA 0x18F5B; call    sub_1E3D7 */
m->ip=0x8f5b;
su_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x8f5e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18F5E: ; /* IDA 0x18F5E; mov     al, 21h ; '!' */
m->ip=0x8f5e;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(33));
goto L_18F60;
L_18F60: ; /* IDA 0x18F60; mov     dx, 148h */
m->ip=0x8f60;
su_skip_first_enter=0;
m->dx=(uint16_t)(328);
goto L_18F63;
L_18F63: ; /* IDA 0x18F63; mov     bx, 0F8h */
m->ip=0x8f63;
su_skip_first_enter=0;
m->bx=(uint16_t)(248);
goto L_18F66;
L_18F66: ; /* IDA 0x18F66; mov     cx, 206h */
m->ip=0x8f66;
su_skip_first_enter=0;
m->cx=(uint16_t)(518);
goto L_18F69;
L_18F69: ; /* IDA 0x18F69; call    sub_1E3D7 */
m->ip=0x8f69;
su_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x8f6c);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18F6C: ; /* IDA 0x18F6C; retn */
m->ip=0x8f6c;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_18F6D: ; /* IDA 0x18F6D; xor     al, al */
m->ip=0x8f6d;
/* Original entry: sub_18F6D */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x8f6d,h->user);
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_18F6F;
L_18F6F: ; /* IDA 0x18F6F; mov     dx, 0Ah */
m->ip=0x8f6f;
su_skip_first_enter=0;
m->dx=(uint16_t)(10);
goto L_18F72;
L_18F72: ; /* IDA 0x18F72; mov     bx, 5 */
m->ip=0x8f72;
su_skip_first_enter=0;
m->bx=(uint16_t)(5);
goto L_18F75;
L_18F75: ; /* IDA 0x18F75; mov     cx, 0C0Fh */
m->ip=0x8f75;
su_skip_first_enter=0;
m->cx=(uint16_t)(3087);
goto L_18F78;
L_18F78: ; /* IDA 0x18F78; call    sub_1895D */
m->ip=0x8f78;
su_skip_first_enter=0;
{uint16_t target=35165;ec_push(m,0x8f7b);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_18F7B: ; /* IDA 0x18F7B; retn */
m->ip=0x8f7b;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_18F7C: ; /* IDA 0x18F7C; mov     al, [si+3Fh] */
m->ip=0x8f7c;
/* Original entry: sub_18F7C */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x8f7c,h->user);
vg_al(m,(uint8_t)(dl_read8(m,m->ds,(uint16_t)(m->si+63))));
goto L_18F7F;
L_18F7F: ; /* IDA 0x18F7F; mov     [si+2], al */
m->ip=0x8f7f;
su_skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->si+2),(uint8_t)m->ax);
goto L_18F82;
L_18F82: ; /* IDA 0x18F82; mov     byte ptr ds:5221h, 0FFh */
m->ip=0x8f82;
su_skip_first_enter=0;
dl_write8(m,m->ds,21025,255);
goto L_18F87;
L_18F87: ; /* IDA 0x18F87; mov     ax, 0D0A1h */
m->ip=0x8f87;
su_skip_first_enter=0;
m->ax=(uint16_t)(53409);
goto L_18F8A;
L_18F8A: ; /* IDA 0x18F8A; mov     bx, 5222h */
m->ip=0x8f8a;
su_skip_first_enter=0;
m->bx=(uint16_t)(21026);
goto L_18F8D;
L_18F8D: ; /* IDA 0x18F8D; mov     cx, 6 */
m->ip=0x8f8d;
su_skip_first_enter=0;
m->cx=(uint16_t)(6);
goto L_18F90;
L_18F90: ; /* IDA 0x18F90; mov     [bx], ax */
m->ip=0x8f90;
su_skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->bx),m->ax);
goto L_18F92;
L_18F92: ; /* IDA 0x18F92; inc     bx */
m->ip=0x8f92;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_inc(m,m->bx,16));
goto L_18F93;
L_18F93: ; /* IDA 0x18F93; inc     bx */
m->ip=0x8f93;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_inc(m,m->bx,16));
goto L_18F94;
L_18F94: ; /* IDA 0x18F94; loop    loc_18F90 */
m->ip=0x8f94;
su_skip_first_enter=0;
if(--m->cx) {goto L_18F90;}
goto L_18F96;
L_18F96: ; /* IDA 0x18F96; retn */
m->ip=0x8f96;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_18FC9: ; /* IDA 0x18FC9; push    ds */
m->ip=0x8fc9;
/* Original entry: sub_18FC9 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x8fc9,h->user);
ec_push(m,m->ds);
goto L_18FCA;
L_18FCA: ; /* IDA 0x18FCA; push    si */
m->ip=0x8fca;
su_skip_first_enter=0;
ec_push(m,m->si);
goto L_18FCB;
L_18FCB: ; /* IDA 0x18FCB; push    bx */
m->ip=0x8fcb;
su_skip_first_enter=0;
ec_push(m,m->bx);
goto L_18FCC;
L_18FCC: ; /* IDA 0x18FCC; mov     ah, 48h ; 'H' */
m->ip=0x8fcc;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(72));
goto L_18FCE;
L_18FCE: ; /* IDA 0x18FCE; mov     bx, 148h */
m->ip=0x8fce;
su_skip_first_enter=0;
m->bx=(uint16_t)(328);
goto L_18FD1;
L_18FD1: ; /* IDA 0x18FD1; int     21h; DOS - 2+ - ALLOCATE MEMORY */
m->ip=0x8fd1;
su_skip_first_enter=0;
su_interrupt(m,33,0x8fd3,h);
goto L_18FD3;
L_18FD3: ; /* IDA 0x18FD3; mov     cs:word_18FC7, ax */
m->ip=0x8fd3;
su_skip_first_enter=0;
dl_write16(m,m->cs,36807,m->ax);
goto L_18FD7;
L_18FD7: ; /* IDA 0x18FD7; mov     ax, 3D00h */
m->ip=0x8fd7;
su_skip_first_enter=0;
m->ax=(uint16_t)(15616);
goto L_18FDA;
L_18FDA: ; /* IDA 0x18FDA; mov     dx, 8FBBh */
m->ip=0x8fda;
su_skip_first_enter=0;
m->dx=(uint16_t)(36795);
goto L_18FDD;
L_18FDD: ; /* IDA 0x18FDD; int     21h; DOS - 2+ - OPEN DISK FILE WITH HANDLE */
m->ip=0x8fdd;
su_skip_first_enter=0;
su_interrupt(m,33,0x8fdf,h);
goto L_18FDF;
L_18FDF: ; /* IDA 0x18FDF; mov     bx, ax */
m->ip=0x8fdf;
su_skip_first_enter=0;
m->bx=(uint16_t)(m->ax);
goto L_18FE1;
L_18FE1: ; /* IDA 0x18FE1; mov     cx, 147Ah */
m->ip=0x8fe1;
su_skip_first_enter=0;
m->cx=(uint16_t)(5242);
goto L_18FE4;
L_18FE4: ; /* IDA 0x18FE4; mov     ds, cs:word_18FC7 */
m->ip=0x8fe4;
su_skip_first_enter=0;
m->ds=(uint16_t)(dl_read16(m,m->cs,36807));
goto L_18FE9;
L_18FE9: ; /* IDA 0x18FE9; xor     dx, dx */
m->ip=0x8fe9;
su_skip_first_enter=0;
{uint16_t result=((m->dx)^(m->dx))&65535;ec_logic(m,result,16);m->dx=(uint16_t)(result);}
goto L_18FEB;
L_18FEB: ; /* IDA 0x18FEB; mov     ah, 3Fh */
m->ip=0x8feb;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(63));
goto L_18FED;
L_18FED: ; /* IDA 0x18FED; int     21h; DOS - 2+ - READ FROM FILE WITH HANDLE */
m->ip=0x8fed;
su_skip_first_enter=0;
su_interrupt(m,33,0x8fef,h);
goto L_18FEF;
L_18FEF: ; /* IDA 0x18FEF; mov     ah, 3Eh */
m->ip=0x8fef;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(62));
goto L_18FF1;
L_18FF1: ; /* IDA 0x18FF1; int     21h; DOS - 2+ - CLOSE A FILE WITH HANDLE */
m->ip=0x8ff1;
su_skip_first_enter=0;
su_interrupt(m,33,0x8ff3,h);
goto L_18FF3;
L_18FF3: ; /* IDA 0x18FF3; mov     ds, cs:word_10D52 */
m->ip=0x8ff3;
su_skip_first_enter=0;
m->ds=(uint16_t)(dl_read16(m,m->cs,3410));
goto L_18FF8;
L_18FF8: ; /* IDA 0x18FF8; cmp     byte ptr ds:5221h, 0FFh */
m->ip=0x8ff8;
su_skip_first_enter=0;
ec_cmp(m,dl_read8(m,m->ds,21025),255,8);
goto L_18FFD;
L_18FFD: ; /* IDA 0x18FFD; jnz     short loc_19004 */
m->ip=0x8ffd;
su_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_19004;}
goto L_18FFF;
L_18FFF: ; /* IDA 0x18FFF; mov     byte ptr ds:5221h, 91h */
m->ip=0x8fff;
su_skip_first_enter=0;
dl_write8(m,m->ds,21025,145);
goto L_19004;
L_19004: ; /* IDA 0x19004; call    sub_101B4 */
m->ip=0x9004;
su_skip_first_enter=0;
{uint16_t target=436;ec_push(m,0x9007);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19007: ; /* IDA 0x19007; call    sub_1917C */
m->ip=0x9007;
su_skip_first_enter=0;
{uint16_t target=37244;ec_push(m,0x900a);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1900A: ; /* IDA 0x1900A; call    sub_1915B */
m->ip=0x900a;
su_skip_first_enter=0;
{uint16_t target=37211;ec_push(m,0x900d);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1900D: ; /* IDA 0x1900D; xor     si, si */
m->ip=0x900d;
su_skip_first_enter=0;
{uint16_t result=((m->si)^(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_1900F;
L_1900F: ; /* IDA 0x1900F; mov     di, 0 */
m->ip=0x900f;
su_skip_first_enter=0;
m->di=(uint16_t)(0);
goto L_19012;
L_19012: ; /* IDA 0x19012; call    sub_19223 */
m->ip=0x9012;
su_skip_first_enter=0;
{uint16_t target=37411;ec_push(m,0x9015);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19015: ; /* IDA 0x19015; call    sub_1928A */
m->ip=0x9015;
su_skip_first_enter=0;
{uint16_t target=37514;ec_push(m,0x9018);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19018: ; /* IDA 0x19018; call    sub_101DB */
m->ip=0x9018;
su_skip_first_enter=0;
{uint16_t target=475;ec_push(m,0x901b);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1901B: ; /* IDA 0x1901B; call    sub_121E7 */
m->ip=0x901b;
su_skip_first_enter=0;
{uint16_t target=8679;ec_push(m,0x901e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1901E: ; /* IDA 0x1901E; jb      short loc_19041 */
m->ip=0x901e;
su_skip_first_enter=0;
if(m->flags&1) {goto L_19041;}
goto L_19020;
L_19020: ; /* IDA 0x19020; push    cx */
m->ip=0x9020;
su_skip_first_enter=0;
ec_push(m,m->cx);
goto L_19021;
L_19021: ; /* IDA 0x19021; push    dx */
m->ip=0x9021;
su_skip_first_enter=0;
ec_push(m,m->dx);
goto L_19022;
L_19022: ; /* IDA 0x19022; call    sub_1E453 */
m->ip=0x9022;
su_skip_first_enter=0;
{uint16_t target=58451;ec_push(m,0x9025);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19025: ; /* IDA 0x19025; pop     dx */
m->ip=0x9025;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_pop(m));
goto L_19026;
L_19026: ; /* IDA 0x19026; pop     cx */
m->ip=0x9026;
su_skip_first_enter=0;
m->cx=(uint16_t)(ec_pop(m));
goto L_19027;
L_19027: ; /* IDA 0x19027; sub     al, 20h ; ' ' */
m->ip=0x9027;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(ec_math(m,(uint8_t)m->ax,32,8,0,1)));
goto L_19029;
L_19029: ; /* IDA 0x19029; jb      short loc_1901B */
m->ip=0x9029;
su_skip_first_enter=0;
if(m->flags&1) {goto L_1901B;}
goto L_1902B;
L_1902B: ; /* IDA 0x1902B; jz      short loc_19046 */
m->ip=0x902b;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_19046;}
goto L_1902D;
L_1902D: ; /* IDA 0x1902D; mov     bl, al */
m->ip=0x902d;
su_skip_first_enter=0;
gw_bl(m,(uint8_t)((uint8_t)m->ax));
goto L_1902F;
L_1902F: ; /* IDA 0x1902F; xor     bh, bh */
m->ip=0x902f;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->bx>>8))^((uint8_t)(m->bx>>8)))&255;ec_logic(m,result,8);gw_bh(m,(uint8_t)(result));}
goto L_19031;
L_19031: ; /* IDA 0x19031; dec     bx */
m->ip=0x9031;
su_skip_first_enter=0;
m->bx=(uint16_t)(st_dec(m,m->bx,16));
goto L_19032;
L_19032: ; /* IDA 0x19032; shl     bx, 1 */
m->ip=0x9032;
su_skip_first_enter=0;
m->bx=(uint16_t)(dl_shift(m,m->bx,1,16,4));
goto L_19034;
L_19034: ; /* IDA 0x19034; call    sub_101B4 */
m->ip=0x9034;
su_skip_first_enter=0;
{uint16_t target=436;ec_push(m,0x9037);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19037: ; /* IDA 0x19037; call    cs:funcs_19037[bx] */
m->ip=0x9037;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 6
{uint16_t target=dl_read16(m,m->cs,(uint16_t)(m->bx+36761));ec_push(m,0x903c);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
#else
{uint16_t target=dl_read16(m,m->cs,(uint16_t)(m->bx+36759));ec_push(m,0x903c);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
#endif
L_1903C: ; /* IDA 0x1903C; call    sub_101DB */
m->ip=0x903c;
su_skip_first_enter=0;
{uint16_t target=475;ec_push(m,0x903f);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1903F: ; /* IDA 0x1903F; jmp     short loc_1901B */
m->ip=0x903f;
su_skip_first_enter=0;
goto L_1901B;
L_19041: ; /* IDA 0x19041; call    sub_18F7C */
m->ip=0x9041;
su_skip_first_enter=0;
{uint16_t target=36732;ec_push(m,0x9044);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19044: ; /* IDA 0x19044; jmp     short loc_1904C */
m->ip=0x9044;
su_skip_first_enter=0;
goto L_1904C;
L_19046: ; /* IDA 0x19046; pop     bx */
m->ip=0x9046;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_19047;
L_19047: ; /* IDA 0x19047; mov     byte ptr [bx+2], 7Fh */
m->ip=0x9047;
su_skip_first_enter=0;
dl_write8(m,m->ds,(uint16_t)(m->bx+2),127);
goto L_1904B;
L_1904B: ; /* IDA 0x1904B; push    bx */
m->ip=0x904b;
su_skip_first_enter=0;
ec_push(m,m->bx);
goto L_1904C;
L_1904C: ; /* IDA 0x1904C; call    sub_19214 */
m->ip=0x904c;
su_skip_first_enter=0;
{uint16_t target=37396;ec_push(m,0x904f);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1904F: ; /* IDA 0x1904F; push    es */
m->ip=0x904f;
su_skip_first_enter=0;
ec_push(m,m->es);
goto L_19050;
L_19050: ; /* IDA 0x19050; mov     es, cs:word_18FC7 */
m->ip=0x9050;
su_skip_first_enter=0;
m->es=(uint16_t)(dl_read16(m,m->cs,36807));
goto L_19055;
L_19055: ; /* IDA 0x19055; mov     ah, 49h */
m->ip=0x9055;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(73));
goto L_19057;
L_19057: ; /* IDA 0x19057; int     21h; DOS - 2+ - FREE MEMORY */
m->ip=0x9057;
su_skip_first_enter=0;
su_interrupt(m,33,0x9059,h);
goto L_19059;
L_19059: ; /* IDA 0x19059; pop     es */
m->ip=0x9059;
su_skip_first_enter=0;
m->es=(uint16_t)(ec_pop(m));
goto L_1905A;
L_1905A: ; /* IDA 0x1905A; pop     bx */
m->ip=0x905a;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_1905B;
L_1905B: ; /* IDA 0x1905B; pop     si */
m->ip=0x905b;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_pop(m));
goto L_1905C;
L_1905C: ; /* IDA 0x1905C; pop     ds */
m->ip=0x905c;
su_skip_first_enter=0;
m->ds=(uint16_t)(ec_pop(m));
goto L_1905D;
L_1905D: ; /* IDA 0x1905D; retn */
m->ip=0x905d;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1905E: ; /* IDA 0x1905E; mov     word ptr [si+5222h], 40A1h */
m->ip=0x905e;
/* Original entry: sub_1905E */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x905e,h->user);
dl_write16(m,m->ds,(uint16_t)(m->si+21026),16545);
goto L_19064;
L_19064: ; /* IDA 0x19064; and     si, si */
m->ip=0x9064;
su_skip_first_enter=0;
{uint16_t result=((m->si)&(m->si))&65535;ec_logic(m,result,16);m->si=(uint16_t)(result);}
goto L_19066;
L_19066: ; /* IDA 0x19066; jz      short loc_1906A */
m->ip=0x9066;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_1906A;}
goto L_19068;
L_19068: ; /* IDA 0x19068; dec     si */
m->ip=0x9068;
su_skip_first_enter=0;
m->si=(uint16_t)(st_dec(m,m->si,16));
goto L_19069;
L_19069: ; /* IDA 0x19069; dec     si */
m->ip=0x9069;
su_skip_first_enter=0;
m->si=(uint16_t)(st_dec(m,m->si,16));
goto L_1906A;
L_1906A: ; /* IDA 0x1906A; call    sub_19223 */
m->ip=0x906a;
su_skip_first_enter=0;
{uint16_t target=37411;ec_push(m,0x906d);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1906D: ; /* IDA 0x1906D; retn */
m->ip=0x906d;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1906E: ; /* IDA 0x1906E; cmp     si, 0Ah */
m->ip=0x906e;
/* Original entry: sub_1906E */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x906e,h->user);
ec_cmp(m,m->si,10,16);
goto L_19071;
L_19071: ; /* IDA 0x19071; jz      short loc_19075 */
m->ip=0x9071;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_19075;}
goto L_19073;
L_19073: ; /* IDA 0x19073; inc     si */
m->ip=0x9073;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_inc(m,m->si,16));
goto L_19074;
L_19074: ; /* IDA 0x19074; inc     si */
m->ip=0x9074;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_inc(m,m->si,16));
goto L_19075;
L_19075: ; /* IDA 0x19075; mov     word ptr [si+5222h], 40A1h */
m->ip=0x9075;
su_skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+21026),16545);
goto L_1907B;
L_1907B: ; /* IDA 0x1907B; call    sub_19223 */
m->ip=0x907b;
su_skip_first_enter=0;
{uint16_t target=37411;ec_push(m,0x907e);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1907E: ; /* IDA 0x1907E; retn */
m->ip=0x907e;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1907F: ; /* IDA 0x1907F; mov     ax, di */
m->ip=0x907f;
/* Original entry: sub_1907F */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x907f,h->user);
m->ax=(uint16_t)(m->di);
goto L_19081;
L_19081: ; /* IDA 0x19081; mov     cx, 0C0h */
m->ip=0x9081;
su_skip_first_enter=0;
m->cx=(uint16_t)(192);
goto L_19084;
L_19084: ; /* IDA 0x19084; call    sub_192EE */
m->ip=0x9084;
su_skip_first_enter=0;
{uint16_t target=37614;ec_push(m,0x9087);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19087: ; /* IDA 0x19087; mov     di, ax */
m->ip=0x9087;
su_skip_first_enter=0;
m->di=(uint16_t)(m->ax);
goto L_19089;
L_19089: ; /* IDA 0x19089; call    sub_1928A */
m->ip=0x9089;
su_skip_first_enter=0;
{uint16_t target=37514;ec_push(m,0x908c);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1908C: ; /* IDA 0x1908C; retn */
m->ip=0x908c;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1908D: ; /* IDA 0x1908D; mov     ax, di */
m->ip=0x908d;
/* Original entry: sub_1908D */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x908d,h->user);
m->ax=(uint16_t)(m->di);
goto L_1908F;
L_1908F: ; /* IDA 0x1908F; mov     cx, 0C0h */
m->ip=0x908f;
su_skip_first_enter=0;
m->cx=(uint16_t)(192);
goto L_19092;
L_19092: ; /* IDA 0x19092; call    sub_192E3 */
m->ip=0x9092;
su_skip_first_enter=0;
{uint16_t target=37603;ec_push(m,0x9095);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19095: ; /* IDA 0x19095; cmp     ax, 13BAh */
m->ip=0x9095;
su_skip_first_enter=0;
ec_cmp(m,m->ax,5050,16);
goto L_19098;
L_19098: ; /* IDA 0x19098; jbe     short loc_1909D */
m->ip=0x9098;
su_skip_first_enter=0;
if(m->flags&(1|0x40)) {goto L_1909D;}
goto L_1909A;
L_1909A: ; /* IDA 0x1909A; mov     ax, 13BAh */
m->ip=0x909a;
su_skip_first_enter=0;
m->ax=(uint16_t)(5050);
goto L_1909D;
L_1909D: ; /* IDA 0x1909D; mov     di, ax */
m->ip=0x909d;
su_skip_first_enter=0;
m->di=(uint16_t)(m->ax);
goto L_1909F;
L_1909F: ; /* IDA 0x1909F; call    sub_1928A */
m->ip=0x909f;
su_skip_first_enter=0;
{uint16_t target=37514;ec_push(m,0x90a2);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_190A2: ; /* IDA 0x190A2; retn */
m->ip=0x90a2;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_190A3: ; /* IDA 0x190A3; sub     cx, 0D8h */
m->ip=0x90a3;
/* Original entry: sub_190A3 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x90a3,h->user);
m->cx=(uint16_t)(ec_math(m,m->cx,216,16,0,1));
goto L_190A7;
L_190A7: ; /* IDA 0x190A7; shr     cx, 1 */
m->ip=0x90a7;
su_skip_first_enter=0;
m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
goto L_190A9;
L_190A9: ; /* IDA 0x190A9; shr     cx, 1 */
m->ip=0x90a9;
su_skip_first_enter=0;
m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
goto L_190AB;
L_190AB: ; /* IDA 0x190AB; shr     cx, 1 */
m->ip=0x90ab;
su_skip_first_enter=0;
m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
goto L_190AD;
L_190AD: ; /* IDA 0x190AD; shr     cx, 1 */
m->ip=0x90ad;
su_skip_first_enter=0;
m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
goto L_190AF;
L_190AF: ; /* IDA 0x190AF; shr     cx, 1 */
m->ip=0x90af;
su_skip_first_enter=0;
m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,5));
goto L_190B1;
L_190B1: ; /* IDA 0x190B1; jb      short locret_190BF */
m->ip=0x90b1;
su_skip_first_enter=0;
if(m->flags&1) {goto L_190BF;}
goto L_190B3;
L_190B3: ; /* IDA 0x190B3; shl     cx, 1 */
m->ip=0x90b3;
su_skip_first_enter=0;
m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
goto L_190B5;
L_190B5: ; /* IDA 0x190B5; mov     bx, cx */
m->ip=0x90b5;
su_skip_first_enter=0;
m->bx=(uint16_t)(m->cx);
goto L_190B7;
L_190B7: ; /* IDA 0x190B7; mov     di, cs:[bx-7059h] */
m->ip=0x90b7;
su_skip_first_enter=0;
m->di=(uint16_t)(dl_read16(m,m->cs,(uint16_t)(m->bx+36775)));
goto L_190BC;
L_190BC: ; /* IDA 0x190BC; call    sub_1928A */
m->ip=0x90bc;
su_skip_first_enter=0;
{uint16_t target=37514;ec_push(m,0x90bf);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_190BF: ; /* IDA 0x190BF; retn */
m->ip=0x90bf;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_190C0: ; /* IDA 0x190C0; sub     cx, 0D2h */
m->ip=0x90c0;
/* Original entry: sub_190C0 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x90c0,h->user);
m->cx=(uint16_t)(ec_math(m,m->cx,210,16,0,1));
goto L_190C4;
L_190C4: ; /* IDA 0x190C4; jb      short locret_190D7 */
m->ip=0x90c4;
su_skip_first_enter=0;
if(m->flags&1) {goto L_190D7;}
goto L_190C6;
L_190C6: ; /* IDA 0x190C6; cmp     cx, 140h */
m->ip=0x90c6;
su_skip_first_enter=0;
ec_cmp(m,m->cx,320,16);
goto L_190CA;
L_190CA: ; /* IDA 0x190CA; jnb     short locret_190D7 */
m->ip=0x90ca;
su_skip_first_enter=0;
if(!(m->flags&1)) {goto L_190D7;}
goto L_190CC;
L_190CC: ; /* IDA 0x190CC; sub     dx, 0EEh */
m->ip=0x90cc;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_math(m,m->dx,238,16,0,1));
goto L_190D0;
L_190D0: ; /* IDA 0x190D0; jb      short locret_190D7 */
m->ip=0x90d0;
su_skip_first_enter=0;
if(m->flags&1) {goto L_190D7;}
goto L_190D2;
L_190D2: ; /* IDA 0x190D2; cmp     dx, 78h ; 'x' */
m->ip=0x90d2;
su_skip_first_enter=0;
ec_cmp(m,m->dx,120,16);
goto L_190D5;
L_190D5: ; /* IDA 0x190D5; jb      short loc_190D8 */
m->ip=0x90d5;
su_skip_first_enter=0;
if(m->flags&1) {goto L_190D8;}
goto L_190D7;
L_190D7: ; /* IDA 0x190D7; retn */
m->ip=0x90d7;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_190D8: ; /* IDA 0x190D8; mov     ax, cx */
m->ip=0x90d8;
su_skip_first_enter=0;
m->ax=(uint16_t)(m->cx);
goto L_190DA;
L_190DA: ; /* IDA 0x190DA; mov     bl, 14h */
m->ip=0x90da;
su_skip_first_enter=0;
gw_bl(m,(uint8_t)(20));
goto L_190DC;
L_190DC: ; /* IDA 0x190DC; div     bl */
m->ip=0x90dc;
su_skip_first_enter=0;
dl_div(m,(uint8_t)m->bx,8);
goto L_190DE;
L_190DE: ; /* IDA 0x190DE; cmp     ah, 10h */
m->ip=0x90de;
su_skip_first_enter=0;
ec_cmp(m,(uint8_t)(m->ax>>8),16,8);
goto L_190E1;
L_190E1: ; /* IDA 0x190E1; jnb     short locret_190D7 */
m->ip=0x90e1;
su_skip_first_enter=0;
if(!(m->flags&1)) {goto L_190D7;}
goto L_190E3;
L_190E3: ; /* IDA 0x190E3; mov     cx, ax */
m->ip=0x90e3;
su_skip_first_enter=0;
m->cx=(uint16_t)(m->ax);
goto L_190E5;
L_190E5: ; /* IDA 0x190E5; mov     ax, dx */
m->ip=0x90e5;
su_skip_first_enter=0;
m->ax=(uint16_t)(m->dx);
goto L_190E7;
L_190E7: ; /* IDA 0x190E7; mov     bl, 14h */
m->ip=0x90e7;
su_skip_first_enter=0;
gw_bl(m,(uint8_t)(20));
goto L_190E9;
L_190E9: ; /* IDA 0x190E9; div     bl */
m->ip=0x90e9;
su_skip_first_enter=0;
dl_div(m,(uint8_t)m->bx,8);
goto L_190EB;
L_190EB: ; /* IDA 0x190EB; cmp     ah, 10h */
m->ip=0x90eb;
su_skip_first_enter=0;
ec_cmp(m,(uint8_t)(m->ax>>8),16,8);
goto L_190EE;
L_190EE: ; /* IDA 0x190EE; jnb     short locret_190D7 */
m->ip=0x90ee;
su_skip_first_enter=0;
if(!(m->flags&1)) {goto L_190D7;}
goto L_190F0;
L_190F0: ; /* IDA 0x190F0; xor     ah, ah */
m->ip=0x90f0;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)(m->ax>>8))^((uint8_t)(m->ax>>8)))&255;ec_logic(m,result,8);vg_ah(m,(uint8_t)(result));}
goto L_190F2;
L_190F2: ; /* IDA 0x190F2; mov     ch, ah */
m->ip=0x90f2;
su_skip_first_enter=0;
vg_ch(m,(uint8_t)((uint8_t)(m->ax>>8)));
goto L_190F4;
L_190F4: ; /* IDA 0x190F4; shl     ax, 1 */
m->ip=0x90f4;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
goto L_190F6;
L_190F6: ; /* IDA 0x190F6; shl     ax, 1 */
m->ip=0x90f6;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
goto L_190F8;
L_190F8: ; /* IDA 0x190F8; shl     ax, 1 */
m->ip=0x90f8;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
goto L_190FA;
L_190FA: ; /* IDA 0x190FA; shl     ax, 1 */
m->ip=0x90fa;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_shift(m,m->ax,1,16,4));
goto L_190FC;
L_190FC: ; /* IDA 0x190FC; add     cx, ax */
m->ip=0x90fc;
su_skip_first_enter=0;
m->cx=(uint16_t)(ec_math(m,m->cx,m->ax,16,0,0));
goto L_190FE;
L_190FE: ; /* IDA 0x190FE; mov     ax, di */
m->ip=0x90fe;
su_skip_first_enter=0;
m->ax=(uint16_t)(m->di);
goto L_19100;
L_19100: ; /* IDA 0x19100; and     cx, cx */
m->ip=0x9100;
su_skip_first_enter=0;
{uint16_t result=((m->cx)&(m->cx))&65535;ec_logic(m,result,16);m->cx=(uint16_t)(result);}
goto L_19102;
L_19102: ; /* IDA 0x19102; jz      short loc_19109 */
m->ip=0x9102;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_19109;}
goto L_19104;
L_19104: ; /* IDA 0x19104; shl     cx, 1 */
m->ip=0x9104;
su_skip_first_enter=0;
m->cx=(uint16_t)(dl_shift(m,m->cx,1,16,4));
goto L_19106;
L_19106: ; /* IDA 0x19106; call    sub_192E3 */
m->ip=0x9106;
su_skip_first_enter=0;
{uint16_t target=37603;ec_push(m,0x9109);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19109: ; /* IDA 0x19109; push    es */
m->ip=0x9109;
su_skip_first_enter=0;
ec_push(m,m->es);
goto L_1910A;
L_1910A: ; /* IDA 0x1910A; push    di */
m->ip=0x910a;
su_skip_first_enter=0;
ec_push(m,m->di);
goto L_1910B;
L_1910B: ; /* IDA 0x1910B; mov     es, cs:word_18FC7 */
m->ip=0x910b;
su_skip_first_enter=0;
m->es=(uint16_t)(dl_read16(m,m->cs,36807));
goto L_19110;
L_19110: ; /* IDA 0x19110; mov     di, ax */
m->ip=0x9110;
su_skip_first_enter=0;
m->di=(uint16_t)(m->ax);
goto L_19112;
L_19112: ; /* IDA 0x19112; mov     ax, es:[di] */
m->ip=0x9112;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->di)));
goto L_19115;
L_19115: ; /* IDA 0x19115; xchg    ah, al */
m->ip=0x9115;
su_skip_first_enter=0;
{uint16_t a=(uint8_t)(m->ax>>8),b=(uint8_t)m->ax;vg_ah(m,(uint8_t)(b));vg_al(m,(uint8_t)(a));}
goto L_19117;
L_19117: ; /* IDA 0x19117; sub     ax, 1000h */
m->ip=0x9117;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_math(m,m->ax,4096,16,0,1));
goto L_1911A;
L_1911A: ; /* IDA 0x1911A; pop     di */
m->ip=0x911a;
su_skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_1911B;
L_1911B: ; /* IDA 0x1911B; pop     es */
m->ip=0x911b;
su_skip_first_enter=0;
m->es=(uint16_t)(ec_pop(m));
goto L_1911C;
L_1911C: ; /* IDA 0x1911C; xchg    al, ah */
m->ip=0x911c;
su_skip_first_enter=0;
{uint16_t a=(uint8_t)m->ax,b=(uint8_t)(m->ax>>8);vg_al(m,(uint8_t)(b));vg_ah(m,(uint8_t)(a));}
goto L_1911E;
L_1911E: ; /* IDA 0x1911E; mov     [si+5222h], ax */
m->ip=0x911e;
su_skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si+21026),m->ax);
goto L_19122;
L_19122: ; /* IDA 0x19122; cmp     si, 0Ah */
m->ip=0x9122;
su_skip_first_enter=0;
ec_cmp(m,m->si,10,16);
goto L_19125;
L_19125: ; /* IDA 0x19125; jz      short loc_19129 */
m->ip=0x9125;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_19129;}
goto L_19127;
L_19127: ; /* IDA 0x19127; inc     si */
m->ip=0x9127;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_inc(m,m->si,16));
goto L_19128;
L_19128: ; /* IDA 0x19128; inc     si */
m->ip=0x9128;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_inc(m,m->si,16));
goto L_19129;
L_19129: ; /* IDA 0x19129; call    sub_19223 */
m->ip=0x9129;
su_skip_first_enter=0;
{uint16_t target=37411;ec_push(m,0x912c);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1912C: ; /* IDA 0x1912C; retn */
m->ip=0x912c;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1912D: ; /* IDA 0x1912D; mov     al, ds:5221h */
m->ip=0x912d;
/* Original entry: sub_1912D */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x912d,h->user);
vg_al(m,(uint8_t)(dl_read8(m,m->ds,21025)));
goto L_19130;
L_19130: ; /* IDA 0x19130; dec     al */
m->ip=0x9130;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(st_dec(m,(uint8_t)m->ax,8)));
goto L_19132;
L_19132: ; /* IDA 0x19132; cmp     al, 93h */
m->ip=0x9132;
su_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,147,8);
goto L_19134;
L_19134: ; /* IDA 0x19134; jb      short loc_19138 */
m->ip=0x9134;
su_skip_first_enter=0;
if(m->flags&1) {goto L_19138;}
goto L_19136;
L_19136: ; /* IDA 0x19136; mov     al, 92h */
m->ip=0x9136;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 8
vg_al(m,(uint8_t)(147));
#else
vg_al(m,(uint8_t)(146));
#endif
goto L_19138;
L_19138: ; /* IDA 0x19138; call    sub_192F5 */
m->ip=0x9138;
su_skip_first_enter=0;
{uint16_t target=37621;ec_push(m,0x913b);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1913B: ; /* IDA 0x1913B; jb      short loc_19130 */
m->ip=0x913b;
su_skip_first_enter=0;
if(m->flags&1) {goto L_19130;}
goto L_1913D;
L_1913D: ; /* IDA 0x1913D; mov     ds:5221h, al */
m->ip=0x913d;
su_skip_first_enter=0;
dl_write8(m,m->ds,21025,(uint8_t)m->ax);
goto L_19140;
L_19140: ; /* IDA 0x19140; call    sub_1915B */
m->ip=0x9140;
su_skip_first_enter=0;
{uint16_t target=37211;ec_push(m,0x9143);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19143: ; /* IDA 0x19143; retn */
m->ip=0x9143;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_19144: ; /* IDA 0x19144; mov     al, ds:5221h */
m->ip=0x9144;
/* Original entry: sub_19144 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x9144,h->user);
vg_al(m,(uint8_t)(dl_read8(m,m->ds,21025)));
goto L_19147;
L_19147: ; /* IDA 0x19147; inc     al */
m->ip=0x9147;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
goto L_19149;
L_19149: ; /* IDA 0x19149; cmp     al, 93h */
m->ip=0x9149;
su_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,147,8);
goto L_1914B;
L_1914B: ; /* IDA 0x1914B; jb      short loc_1914F */
m->ip=0x914b;
su_skip_first_enter=0;
if(m->flags&1) {goto L_1914F;}
goto L_1914D;
L_1914D: ; /* IDA 0x1914D; xor     al, al */
m->ip=0x914d;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_1914F;
L_1914F: ; /* IDA 0x1914F; call    sub_192F5 */
m->ip=0x914f;
su_skip_first_enter=0;
{uint16_t target=37621;ec_push(m,0x9152);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19152: ; /* IDA 0x19152; jb      short loc_19147 */
m->ip=0x9152;
su_skip_first_enter=0;
if(m->flags&1) {goto L_19147;}
goto L_19154;
L_19154: ; /* IDA 0x19154; mov     ds:5221h, al */
m->ip=0x9154;
su_skip_first_enter=0;
dl_write8(m,m->ds,21025,(uint8_t)m->ax);
goto L_19157;
L_19157: ; /* IDA 0x19157; call    sub_1915B */
m->ip=0x9157;
su_skip_first_enter=0;
{uint16_t target=37211;ec_push(m,0x915a);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1915A: ; /* IDA 0x1915A; retn */
m->ip=0x915a;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1915B: ; /* IDA 0x1915B; push    si */
m->ip=0x915b;
/* Original entry: sub_1915B */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x915b,h->user);
ec_push(m,m->si);
goto L_1915C;
L_1915C: ; /* IDA 0x1915C; mov     al, ds:5221h */
m->ip=0x915c;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(dl_read8(m,m->ds,21025)));
goto L_1915F;
L_1915F: ; /* IDA 0x1915F; mov     bx, 2D1Ah */
m->ip=0x915f;
su_skip_first_enter=0;
m->bx=(uint16_t)(11546);
goto L_19162;
L_19162: ; /* IDA 0x19162; call    sub_107D2 */
m->ip=0x9162;
su_skip_first_enter=0;
{uint16_t target=2002;ec_push(m,0x9165);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19165: ; /* IDA 0x19165; mov     cx, 6 */
m->ip=0x9165;
su_skip_first_enter=0;
m->cx=(uint16_t)(6);
goto L_19168;
L_19168: ; /* IDA 0x19168; mov     si, 5222h */
m->ip=0x9168;
su_skip_first_enter=0;
m->si=(uint16_t)(21026);
goto L_1916B;
L_1916B: ; /* IDA 0x1916B; cmp     word ptr [si], 0D0A1h */
m->ip=0x916b;
su_skip_first_enter=0;
ec_cmp(m,dl_read16(m,m->ds,(uint16_t)(m->si)),53409,16);
goto L_1916F;
L_1916F: ; /* IDA 0x1916F; jnz     short loc_1917A */
m->ip=0x916f;
su_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_1917A;}
goto L_19171;
L_19171: ; /* IDA 0x19171; mov     ax, 40A1h */
m->ip=0x9171;
su_skip_first_enter=0;
m->ax=(uint16_t)(16545);
goto L_19174;
L_19174: ; /* IDA 0x19174; mov     [si], ax */
m->ip=0x9174;
su_skip_first_enter=0;
dl_write16(m,m->ds,(uint16_t)(m->si),m->ax);
goto L_19176;
L_19176: ; /* IDA 0x19176; inc     si */
m->ip=0x9176;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_inc(m,m->si,16));
goto L_19177;
L_19177: ; /* IDA 0x19177; inc     si */
m->ip=0x9177;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_inc(m,m->si,16));
goto L_19178;
L_19178: ; /* IDA 0x19178; loop    loc_19174 */
m->ip=0x9178;
su_skip_first_enter=0;
if(--m->cx) {goto L_19174;}
goto L_1917A;
L_1917A: ; /* IDA 0x1917A; pop     si */
m->ip=0x917a;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_pop(m));
goto L_1917B;
L_1917B: ; /* IDA 0x1917B; retn */
m->ip=0x917b;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1917C: ; /* IDA 0x1917C; mov     al, 1Fh */
m->ip=0x917c;
/* Original entry: sub_1917C */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x917c,h->user);
vg_al(m,(uint8_t)(31));
goto L_1917E;
L_1917E: ; /* IDA 0x1917E; mov     dx, 0Ch */
m->ip=0x917e;
su_skip_first_enter=0;
m->dx=(uint16_t)(12);
goto L_19181;
L_19181: ; /* IDA 0x19181; mov     bx, 6 */
m->ip=0x9181;
su_skip_first_enter=0;
m->bx=(uint16_t)(6);
goto L_19184;
L_19184: ; /* IDA 0x19184; mov     cx, 1016h */
m->ip=0x9184;
su_skip_first_enter=0;
m->cx=(uint16_t)(4118);
goto L_19187;
L_19187: ; /* IDA 0x19187; call    sub_1895D */
m->ip=0x9187;
su_skip_first_enter=0;
{uint16_t target=35165;ec_push(m,0x918a);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1918A: ; /* IDA 0x1918A; mov     al, 9 */
m->ip=0x918a;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(9));
goto L_1918C;
L_1918C: ; /* IDA 0x1918C; mov     dx, 0C8h */
m->ip=0x918c;
su_skip_first_enter=0;
m->dx=(uint16_t)(200);
goto L_1918F;
L_1918F: ; /* IDA 0x1918F; mov     bx, 88h */
m->ip=0x918f;
su_skip_first_enter=0;
m->bx=(uint16_t)(136);
goto L_19192;
L_19192: ; /* IDA 0x19192; call    sub_10337 */
m->ip=0x9192;
su_skip_first_enter=0;
{uint16_t target=823;ec_push(m,0x9195);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19195: ; /* IDA 0x19195; mov     al, 20h ; ' ' */
m->ip=0x9195;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(32));
goto L_19197;
L_19197: ; /* IDA 0x19197; mov     dx, 1D0h */
m->ip=0x9197;
su_skip_first_enter=0;
m->dx=(uint16_t)(464);
goto L_1919A;
L_1919A: ; /* IDA 0x1919A; mov     bx, 0C0h */
m->ip=0x919a;
su_skip_first_enter=0;
m->bx=(uint16_t)(192);
goto L_1919D;
L_1919D: ; /* IDA 0x1919D; mov     cx, 208h */
m->ip=0x919d;
su_skip_first_enter=0;
m->cx=(uint16_t)(520);
goto L_191A0;
L_191A0: ; /* IDA 0x191A0; call    sub_1E3D7 */
m->ip=0x91a0;
su_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x91a3);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_191A3: ; /* IDA 0x191A3; mov     al, 21h ; '!' */
m->ip=0x91a3;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(33));
goto L_191A5;
L_191A5: ; /* IDA 0x191A5; mov     dx, 160h */
m->ip=0x91a5;
su_skip_first_enter=0;
m->dx=(uint16_t)(352);
goto L_191A8;
L_191A8: ; /* IDA 0x191A8; mov     bx, 0C0h */
m->ip=0x91a8;
su_skip_first_enter=0;
m->bx=(uint16_t)(192);
goto L_191AB;
L_191AB: ; /* IDA 0x191AB; mov     cx, 206h */
m->ip=0x91ab;
su_skip_first_enter=0;
m->cx=(uint16_t)(518);
goto L_191AE;
L_191AE: ; /* IDA 0x191AE; call    sub_1E3D7 */
m->ip=0x91ae;
su_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x91b1);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_191B1: ; /* IDA 0x191B1; mov     al, 22h ; '"' */
m->ip=0x91b1;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(34));
goto L_191B3;
L_191B3: ; /* IDA 0x191B3; mov     dx, 198h */
m->ip=0x91b3;
su_skip_first_enter=0;
m->dx=(uint16_t)(408);
goto L_191B6;
L_191B6: ; /* IDA 0x191B6; mov     bx, 0C0h */
m->ip=0x91b6;
su_skip_first_enter=0;
m->bx=(uint16_t)(192);
goto L_191B9;
L_191B9: ; /* IDA 0x191B9; mov     cx, 206h */
m->ip=0x91b9;
su_skip_first_enter=0;
m->cx=(uint16_t)(518);
goto L_191BC;
L_191BC: ; /* IDA 0x191BC; call    sub_1E3D7 */
m->ip=0x91bc;
su_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x91bf);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_191BF: ; /* IDA 0x191BF; mov     al, 23h ; '#' */
m->ip=0x91bf;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(35));
goto L_191C1;
L_191C1: ; /* IDA 0x191C1; mov     dx, 0C8h */
m->ip=0x91c1;
su_skip_first_enter=0;
m->dx=(uint16_t)(200);
goto L_191C4;
L_191C4: ; /* IDA 0x191C4; mov     bx, 168h */
m->ip=0x91c4;
su_skip_first_enter=0;
m->bx=(uint16_t)(360);
goto L_191C7;
L_191C7: ; /* IDA 0x191C7; mov     cx, 215h */
m->ip=0x91c7;
su_skip_first_enter=0;
m->cx=(uint16_t)(533);
goto L_191CA;
L_191CA: ; /* IDA 0x191CA; call    sub_1E3D7 */
m->ip=0x91ca;
su_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x91cd);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_191CD: ; /* IDA 0x191CD; mov     al, 24h ; '$' */
m->ip=0x91cd;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(36));
goto L_191CF;
L_191CF: ; /* IDA 0x191CF; mov     dx, 170h */
m->ip=0x91cf;
su_skip_first_enter=0;
m->dx=(uint16_t)(368);
goto L_191D2;
L_191D2: ; /* IDA 0x191D2; mov     bx, 168h */
m->ip=0x91d2;
su_skip_first_enter=0;
m->bx=(uint16_t)(360);
goto L_191D5;
L_191D5: ; /* IDA 0x191D5; mov     cx, 215h */
m->ip=0x91d5;
su_skip_first_enter=0;
m->cx=(uint16_t)(533);
goto L_191D8;
L_191D8: ; /* IDA 0x191D8; call    sub_1E3D7 */
m->ip=0x91d8;
su_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x91db);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_191DB: ; /* IDA 0x191DB; mov     al, 25h ; '%' */
m->ip=0x91db;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(37));
goto L_191DD;
L_191DD: ; /* IDA 0x191DD; mov     dx, 0D8h */
m->ip=0x91dd;
su_skip_first_enter=0;
m->dx=(uint16_t)(216);
goto L_191E0;
L_191E0: ; /* IDA 0x191E0; mov     bx, 0D8h */
m->ip=0x91e0;
su_skip_first_enter=0;
m->bx=(uint16_t)(216);
goto L_191E3;
L_191E3: ; /* IDA 0x191E3; mov     cx, 228h */
m->ip=0x91e3;
su_skip_first_enter=0;
m->cx=(uint16_t)(552);
goto L_191E6;
L_191E6: ; /* IDA 0x191E6; call    sub_1E3D7 */
m->ip=0x91e6;
su_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x91e9);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_191E9: ; /* IDA 0x191E9; mov     al, 26h ; '&' */
m->ip=0x91e9;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(38));
goto L_191EB;
L_191EB: ; /* IDA 0x191EB; mov     dx, 0C8h */
m->ip=0x91eb;
su_skip_first_enter=0;
m->dx=(uint16_t)(200);
goto L_191EE;
L_191EE: ; /* IDA 0x191EE; mov     bx, 0E8h */
m->ip=0x91ee;
su_skip_first_enter=0;
m->bx=(uint16_t)(232);
goto L_191F1;
L_191F1: ; /* IDA 0x191F1; mov     cx, 102Ah */
m->ip=0x91f1;
su_skip_first_enter=0;
m->cx=(uint16_t)(4138);
goto L_191F4;
L_191F4: ; /* IDA 0x191F4; call    sub_1E3D7 */
m->ip=0x91f4;
su_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x91f7);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_191F7: ; /* IDA 0x191F7; mov     al, 27h ; ''' */
m->ip=0x91f7;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(39));
goto L_191F9;
L_191F9: ; /* IDA 0x191F9; mov     dx, 110h */
m->ip=0x91f9;
su_skip_first_enter=0;
m->dx=(uint16_t)(272);
goto L_191FC;
L_191FC: ; /* IDA 0x191FC; mov     bx, 90h */
m->ip=0x91fc;
su_skip_first_enter=0;
m->bx=(uint16_t)(144);
goto L_191FF;
L_191FF: ; /* IDA 0x191FF; mov     cx, 407h */
m->ip=0x91ff;
su_skip_first_enter=0;
m->cx=(uint16_t)(1031);
goto L_19202;
L_19202: ; /* IDA 0x19202; call    sub_1E3D7 */
m->ip=0x9202;
su_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x9205);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19205: ; /* IDA 0x19205; mov     al, 28h ; '(' */
m->ip=0x9205;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(40));
goto L_19207;
L_19207: ; /* IDA 0x19207; mov     dx, 110h */
m->ip=0x9207;
su_skip_first_enter=0;
m->dx=(uint16_t)(272);
goto L_1920A;
L_1920A: ; /* IDA 0x1920A; mov     bx, 0B0h */
m->ip=0x920a;
su_skip_first_enter=0;
m->bx=(uint16_t)(176);
goto L_1920D;
L_1920D: ; /* IDA 0x1920D; mov     cx, 407h */
m->ip=0x920d;
su_skip_first_enter=0;
m->cx=(uint16_t)(1031);
goto L_19210;
L_19210: ; /* IDA 0x19210; call    sub_1E3D7 */
m->ip=0x9210;
su_skip_first_enter=0;
{uint16_t target=58327;ec_push(m,0x9213);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19213: ; /* IDA 0x19213; retn */
m->ip=0x9213;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_19214: ; /* IDA 0x19214; xor     al, al */
m->ip=0x9214;
/* Original entry: sub_19214 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x9214,h->user);
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_19216;
L_19216: ; /* IDA 0x19216; mov     dx, 0Ch */
m->ip=0x9216;
su_skip_first_enter=0;
m->dx=(uint16_t)(12);
goto L_19219;
L_19219: ; /* IDA 0x19219; mov     bx, 6 */
m->ip=0x9219;
su_skip_first_enter=0;
m->bx=(uint16_t)(6);
goto L_1921C;
L_1921C: ; /* IDA 0x1921C; mov     cx, 1016h */
m->ip=0x921c;
su_skip_first_enter=0;
m->cx=(uint16_t)(4118);
goto L_1921F;
L_1921F: ; /* IDA 0x1921F; call    sub_1895D */
m->ip=0x921f;
su_skip_first_enter=0;
{uint16_t target=35165;ec_push(m,0x9222);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19222: ; /* IDA 0x19222; retn */
m->ip=0x9222;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_19223: ; /* IDA 0x19223; mov     ax, 2 */
m->ip=0x9223;
/* Original entry: sub_19223 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x9223,h->user);
m->ax=(uint16_t)(2);
goto L_19226;
L_19226: ; /* IDA 0x19226; call    sub_20000 */
m->ip=0x9226;
su_skip_first_enter=0;
mx_far(m,0x9226,0x922b,h);su_invoke_old=0;goto SU_DISPATCH;
L_1922B: ; /* IDA 0x1922B; push    di */
m->ip=0x922b;
su_skip_first_enter=0;
ec_push(m,m->di);
goto L_1922C;
L_1922C: ; /* IDA 0x1922C; push    si */
m->ip=0x922c;
su_skip_first_enter=0;
ec_push(m,m->si);
goto L_1922D;
L_1922D: ; /* IDA 0x1922D; mov     si, 5222h */
m->ip=0x922d;
su_skip_first_enter=0;
m->si=(uint16_t)(21026);
goto L_19230;
L_19230: ; /* IDA 0x19230; mov     ax, 0F01h */
m->ip=0x9230;
su_skip_first_enter=0;
m->ax=(uint16_t)(3841);
goto L_19233;
L_19233: ; /* IDA 0x19233; mov     dx, 160h */
m->ip=0x9233;
su_skip_first_enter=0;
m->dx=(uint16_t)(352);
goto L_19236;
L_19236: ; /* IDA 0x19236; mov     bx, 0A8h */
m->ip=0x9236;
su_skip_first_enter=0;
m->bx=(uint16_t)(168);
goto L_19239;
L_19239: ; /* IDA 0x19239; call    sub_106FD */
m->ip=0x9239;
su_skip_first_enter=0;
{uint16_t target=1789;ec_push(m,0x923c);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1923C: ; /* IDA 0x1923C; mov     si, 5228h */
m->ip=0x923c;
su_skip_first_enter=0;
m->si=(uint16_t)(21032);
goto L_1923F;
L_1923F: ; /* IDA 0x1923F; mov     ax, 0F01h */
m->ip=0x923f;
su_skip_first_enter=0;
m->ax=(uint16_t)(3841);
goto L_19242;
L_19242: ; /* IDA 0x19242; mov     dx, 1C0h */
m->ip=0x9242;
su_skip_first_enter=0;
m->dx=(uint16_t)(448);
goto L_19245;
L_19245: ; /* IDA 0x19245; mov     bx, 0A8h */
m->ip=0x9245;
su_skip_first_enter=0;
m->bx=(uint16_t)(168);
goto L_19248;
L_19248: ; /* IDA 0x19248; call    sub_106FD */
m->ip=0x9248;
su_skip_first_enter=0;
{uint16_t target=1789;ec_push(m,0x924b);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_1924B: ; /* IDA 0x1924B; pop     si */
m->ip=0x924b;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_pop(m));
goto L_1924C;
L_1924C: ; /* IDA 0x1924C; xor     di, di */
m->ip=0x924c;
su_skip_first_enter=0;
{uint16_t result=((m->di)^(m->di))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
goto L_1924E;
L_1924E: ; /* IDA 0x1924E; mov     dx, 160h */
m->ip=0x924e;
su_skip_first_enter=0;
m->dx=(uint16_t)(352);
goto L_19251;
L_19251: ; /* IDA 0x19251; push    dx */
m->ip=0x9251;
su_skip_first_enter=0;
ec_push(m,m->dx);
goto L_19252;
L_19252: ; /* IDA 0x19252; push    si */
m->ip=0x9252;
su_skip_first_enter=0;
ec_push(m,m->si);
goto L_19253;
L_19253: ; /* IDA 0x19253; push    di */
m->ip=0x9253;
su_skip_first_enter=0;
ec_push(m,m->di);
goto L_19254;
L_19254: ; /* IDA 0x19254; cmp     di, 6 */
m->ip=0x9254;
su_skip_first_enter=0;
ec_cmp(m,m->di,6,16);
goto L_19257;
L_19257: ; /* IDA 0x19257; jb      short loc_1925C */
m->ip=0x9257;
su_skip_first_enter=0;
if(m->flags&1) {goto L_1925C;}
goto L_19259;
L_19259: ; /* IDA 0x19259; add     dx, 30h ; '0' */
m->ip=0x9259;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_math(m,m->dx,48,16,0,0));
goto L_1925C;
L_1925C: ; /* IDA 0x1925C; mov     ah, 1 */
m->ip=0x925c;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(1));
goto L_1925E;
L_1925E: ; /* IDA 0x1925E; cmp     si, di */
m->ip=0x925e;
su_skip_first_enter=0;
ec_cmp(m,m->si,m->di,16);
goto L_19260;
L_19260: ; /* IDA 0x19260; jnz     short loc_19264 */
m->ip=0x9260;
su_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_19264;}
goto L_19262;
L_19262: ; /* IDA 0x19262; mov     ah, 0Fh */
m->ip=0x9262;
su_skip_first_enter=0;
vg_ah(m,(uint8_t)(15));
goto L_19264;
L_19264: ; /* IDA 0x19264; mov     si, dx */
m->ip=0x9264;
su_skip_first_enter=0;
m->si=(uint16_t)(m->dx);
goto L_19266;
L_19266: ; /* IDA 0x19266; add     si, 0Eh */
m->ip=0x9266;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_math(m,m->si,14,16,0,0));
goto L_19269;
L_19269: ; /* IDA 0x19269; mov     bx, 0BAh */
m->ip=0x9269;
su_skip_first_enter=0;
m->bx=(uint16_t)(186);
goto L_1926C;
L_1926C: ; /* IDA 0x1926C; mov     di, bx */
m->ip=0x926c;
su_skip_first_enter=0;
m->di=(uint16_t)(m->bx);
goto L_1926E;
L_1926E: ; /* IDA 0x1926E; call    sub_1F1A3 */
m->ip=0x926e;
su_skip_first_enter=0;
{uint16_t target=61859;ec_push(m,0x9271);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_19271: ; /* IDA 0x19271; pop     di */
m->ip=0x9271;
su_skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_19272;
L_19272: ; /* IDA 0x19272; pop     si */
m->ip=0x9272;
su_skip_first_enter=0;
m->si=(uint16_t)(ec_pop(m));
goto L_19273;
L_19273: ; /* IDA 0x19273; pop     dx */
m->ip=0x9273;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_pop(m));
goto L_19274;
L_19274: ; /* IDA 0x19274; add     dx, 10h */
m->ip=0x9274;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_math(m,m->dx,16,16,0,0));
goto L_19277;
L_19277: ; /* IDA 0x19277; inc     di */
m->ip=0x9277;
su_skip_first_enter=0;
m->di=(uint16_t)(ec_inc(m,m->di,16));
goto L_19278;
L_19278: ; /* IDA 0x19278; inc     di */
m->ip=0x9278;
su_skip_first_enter=0;
m->di=(uint16_t)(ec_inc(m,m->di,16));
goto L_19279;
L_19279: ; /* IDA 0x19279; cmp     di, 0Ch */
m->ip=0x9279;
su_skip_first_enter=0;
ec_cmp(m,m->di,12,16);
goto L_1927C;
L_1927C: ; /* IDA 0x1927C; jb      short loc_19251 */
m->ip=0x927c;
su_skip_first_enter=0;
if(m->flags&1) {goto L_19251;}
goto L_1927E;
L_1927E: ; /* IDA 0x1927E; pop     di */
m->ip=0x927e;
su_skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_1927F;
L_1927F: ; /* IDA 0x1927F; mov     ax, 1 */
m->ip=0x927f;
su_skip_first_enter=0;
m->ax=(uint16_t)(1);
goto L_19282;
L_19282: ; /* IDA 0x19282; xor     bx, bx */
m->ip=0x9282;
su_skip_first_enter=0;
{uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
goto L_19284;
L_19284: ; /* IDA 0x19284; call    sub_20000 */
m->ip=0x9284;
su_skip_first_enter=0;
mx_far(m,0x9284,0x9289,h);su_invoke_old=0;goto SU_DISPATCH;
L_19289: ; /* IDA 0x19289; retn */
m->ip=0x9289;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1928A: ; /* IDA 0x1928A; mov     ax, 2 */
m->ip=0x928a;
/* Original entry: sub_1928A */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x928a,h->user);
m->ax=(uint16_t)(2);
goto L_1928D;
L_1928D: ; /* IDA 0x1928D; call    sub_20000 */
m->ip=0x928d;
su_skip_first_enter=0;
mx_far(m,0x928d,0x9292,h);su_invoke_old=0;goto SU_DISPATCH;
L_19292: ; /* IDA 0x19292; push    es */
m->ip=0x9292;
su_skip_first_enter=0;
ec_push(m,m->es);
goto L_19293;
L_19293: ; /* IDA 0x19293; push    di */
m->ip=0x9293;
su_skip_first_enter=0;
ec_push(m,m->di);
goto L_19294;
L_19294: ; /* IDA 0x19294; mov     es, cs:word_18FC7 */
m->ip=0x9294;
su_skip_first_enter=0;
m->es=(uint16_t)(dl_read16(m,m->cs,36807));
goto L_19299;
L_19299: ; /* IDA 0x19299; mov     bx, 0EEh */
m->ip=0x9299;
su_skip_first_enter=0;
m->bx=(uint16_t)(238);
goto L_1929C;
L_1929C: ; /* IDA 0x1929C; mov     ax, di */
m->ip=0x929c;
su_skip_first_enter=0;
m->ax=(uint16_t)(m->di);
goto L_1929E;
L_1929E: ; /* IDA 0x1929E; mov     ch, 6 */
m->ip=0x929e;
su_skip_first_enter=0;
vg_ch(m,(uint8_t)(6));
goto L_192A0;
L_192A0: ; /* IDA 0x192A0; mov     dx, 0D2h */
m->ip=0x92a0;
su_skip_first_enter=0;
m->dx=(uint16_t)(210);
goto L_192A3;
L_192A3: ; /* IDA 0x192A3; mov     cl, 10h */
m->ip=0x92a3;
su_skip_first_enter=0;
vg_cl(m,(uint8_t)(16));
goto L_192A5;
L_192A5: ; /* IDA 0x192A5; push    cx */
m->ip=0x92a5;
su_skip_first_enter=0;
ec_push(m,m->cx);
goto L_192A6;
L_192A6: ; /* IDA 0x192A6; push    ax */
m->ip=0x92a6;
su_skip_first_enter=0;
ec_push(m,m->ax);
goto L_192A7;
L_192A7: ; /* IDA 0x192A7; push    bx */
m->ip=0x92a7;
su_skip_first_enter=0;
ec_push(m,m->bx);
goto L_192A8;
L_192A8: ; /* IDA 0x192A8; push    dx */
m->ip=0x92a8;
su_skip_first_enter=0;
ec_push(m,m->dx);
goto L_192A9;
L_192A9: ; /* IDA 0x192A9; mov     di, ax */
m->ip=0x92a9;
su_skip_first_enter=0;
m->di=(uint16_t)(m->ax);
goto L_192AB;
L_192AB: ; /* IDA 0x192AB; mov     ax, es:[di] */
m->ip=0x92ab;
su_skip_first_enter=0;
m->ax=(uint16_t)(dl_read16(m,m->es,(uint16_t)(m->di)));
goto L_192AE;
L_192AE: ; /* IDA 0x192AE; xchg    ah, al */
m->ip=0x92ae;
su_skip_first_enter=0;
{uint16_t a=(uint8_t)(m->ax>>8),b=(uint8_t)m->ax;vg_ah(m,(uint8_t)(b));vg_al(m,(uint8_t)(a));}
goto L_192B0;
L_192B0: ; /* IDA 0x192B0; sub     ax, 1000h */
m->ip=0x92b0;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_math(m,m->ax,4096,16,0,1));
goto L_192B3;
L_192B3: ; /* IDA 0x192B3; call    sub_1F878 */
m->ip=0x92b3;
su_skip_first_enter=0;
{uint16_t target=63608;ec_push(m,0x92b6);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_192B6: ; /* IDA 0x192B6; mov     cx, ax */
m->ip=0x92b6;
su_skip_first_enter=0;
m->cx=(uint16_t)(m->ax);
goto L_192B8;
L_192B8: ; /* IDA 0x192B8; mov     ax, 9000h */
m->ip=0x92b8;
su_skip_first_enter=0;
m->ax=(uint16_t)(36864);
goto L_192BB;
L_192BB: ; /* IDA 0x192BB; call    loc_1F75E */
m->ip=0x92bb;
su_skip_first_enter=0;
{uint16_t target=63326;ec_push(m,0x92be);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_192BE: ; /* IDA 0x192BE; pop     dx */
m->ip=0x92be;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_pop(m));
goto L_192BF;
L_192BF: ; /* IDA 0x192BF; pop     bx */
m->ip=0x92bf;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_pop(m));
goto L_192C0;
L_192C0: ; /* IDA 0x192C0; pop     ax */
m->ip=0x92c0;
su_skip_first_enter=0;
m->ax=(uint16_t)(ec_pop(m));
goto L_192C1;
L_192C1: ; /* IDA 0x192C1; mov     cx, 2 */
m->ip=0x92c1;
su_skip_first_enter=0;
m->cx=(uint16_t)(2);
goto L_192C4;
L_192C4: ; /* IDA 0x192C4; call    sub_192E3 */
m->ip=0x92c4;
su_skip_first_enter=0;
{uint16_t target=37603;ec_push(m,0x92c7);m->ip=target;su_invoke_old=1;goto SU_DISPATCH;}
L_192C7: ; /* IDA 0x192C7; pop     cx */
m->ip=0x92c7;
su_skip_first_enter=0;
m->cx=(uint16_t)(ec_pop(m));
goto L_192C8;
L_192C8: ; /* IDA 0x192C8; add     dx, 14h */
m->ip=0x92c8;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_math(m,m->dx,20,16,0,0));
goto L_192CB;
L_192CB: ; /* IDA 0x192CB; dec     cl */
m->ip=0x92cb;
su_skip_first_enter=0;
vg_cl(m,(uint8_t)(st_dec(m,(uint8_t)m->cx,8)));
goto L_192CD;
L_192CD: ; /* IDA 0x192CD; jnz     short loc_192A5 */
m->ip=0x92cd;
su_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_192A5;}
goto L_192CF;
L_192CF: ; /* IDA 0x192CF; add     bx, 14h */
m->ip=0x92cf;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,20,16,0,0));
goto L_192D2;
L_192D2: ; /* IDA 0x192D2; dec     ch */
m->ip=0x92d2;
su_skip_first_enter=0;
vg_ch(m,(uint8_t)(st_dec(m,(uint8_t)(m->cx>>8),8)));
goto L_192D4;
L_192D4: ; /* IDA 0x192D4; jnz     short loc_192A0 */
m->ip=0x92d4;
su_skip_first_enter=0;
if(!(m->flags&0x40)) {goto L_192A0;}
goto L_192D6;
L_192D6: ; /* IDA 0x192D6; pop     di */
m->ip=0x92d6;
su_skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_192D7;
L_192D7: ; /* IDA 0x192D7; pop     es */
m->ip=0x92d7;
su_skip_first_enter=0;
m->es=(uint16_t)(ec_pop(m));
goto L_192D8;
L_192D8: ; /* IDA 0x192D8; mov     ax, 1 */
m->ip=0x92d8;
su_skip_first_enter=0;
m->ax=(uint16_t)(1);
goto L_192DB;
L_192DB: ; /* IDA 0x192DB; xor     bx, bx */
m->ip=0x92db;
su_skip_first_enter=0;
{uint16_t result=((m->bx)^(m->bx))&65535;ec_logic(m,result,16);m->bx=(uint16_t)(result);}
goto L_192DD;
L_192DD: ; /* IDA 0x192DD; call    sub_20000 */
m->ip=0x92dd;
su_skip_first_enter=0;
mx_far(m,0x92dd,0x92e2,h);su_invoke_old=0;goto SU_DISPATCH;
L_192E2: ; /* IDA 0x192E2; retn */
m->ip=0x92e2;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_192E3: ; /* IDA 0x192E3; add     ax, cx */
m->ip=0x92e3;
/* Original entry: sub_192E3 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x92e3,h->user);
m->ax=(uint16_t)(ec_math(m,m->ax,m->cx,16,0,0));
goto L_192E5;
L_192E5: ; /* IDA 0x192E5; cmp     ax, 1478h */
m->ip=0x92e5;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 7
ec_cmp(m,m->ax,5241,16);
#else
ec_cmp(m,m->ax,5240,16);
#endif
goto L_192E8;
L_192E8: ; /* IDA 0x192E8; jbe     short locret_192ED */
m->ip=0x92e8;
su_skip_first_enter=0;
if(m->flags&(1|0x40)) {goto L_192ED;}
goto L_192EA;
L_192EA: ; /* IDA 0x192EA; mov     ax, 1478h */
m->ip=0x92ea;
su_skip_first_enter=0;
m->ax=(uint16_t)(5240);
goto L_192ED;
L_192ED: ; /* IDA 0x192ED; retn */
m->ip=0x92ed;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_192EE: ; /* IDA 0x192EE; sub     ax, cx */
m->ip=0x92ee;
/* Original entry: sub_192EE */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x92ee,h->user);
m->ax=(uint16_t)(ec_math(m,m->ax,m->cx,16,0,1));
goto L_192F0;
L_192F0: ; /* IDA 0x192F0; jnb     short locret_192F4 */
m->ip=0x92f0;
su_skip_first_enter=0;
if(!(m->flags&1)) {goto L_192F4;}
goto L_192F2;
L_192F2: ; /* IDA 0x192F2; xor     ax, ax */
m->ip=0x92f2;
su_skip_first_enter=0;
{uint16_t result=((m->ax)^(m->ax))&65535;ec_logic(m,result,16);m->ax=(uint16_t)(result);}
goto L_192F4;
L_192F4: ; /* IDA 0x192F4; retn */
m->ip=0x92f4;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_192F5: ; /* IDA 0x192F5; mov     cx, 7Eh ; '~' */
m->ip=0x92f5;
/* Original entry: sub_192F5 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0x92f5,h->user);
m->cx=(uint16_t)(126);
goto L_192F8;
L_192F8: ; /* IDA 0x192F8; mov     bx, 4240h */
m->ip=0x92f8;
su_skip_first_enter=0;
m->bx=(uint16_t)(16960);
goto L_192FB;
L_192FB: ; /* IDA 0x192FB; cmp     al, [bx+1] */
m->ip=0x92fb;
su_skip_first_enter=0;
ec_cmp(m,(uint8_t)m->ax,dl_read8(m,m->ds,(uint16_t)(m->bx+1)),8);
goto L_192FE;
L_192FE: ; /* IDA 0x192FE; jz      short loc_19307 */
m->ip=0x92fe;
su_skip_first_enter=0;
if(m->flags&0x40) {goto L_19307;}
goto L_19300;
L_19300: ; /* IDA 0x19300; add     bx, 20h ; ' ' */
m->ip=0x9300;
su_skip_first_enter=0;
m->bx=(uint16_t)(ec_math(m,m->bx,32,16,0,0));
goto L_19303;
L_19303: ; /* IDA 0x19303; loop    loc_192FB */
m->ip=0x9303;
su_skip_first_enter=0;
if(--m->cx) {goto L_192FB;}
goto L_19305;
L_19305: ; /* IDA 0x19305; clc */
m->ip=0x9305;
su_skip_first_enter=0;
m->flags&=(uint16_t)~1u;
goto L_19306;
L_19306: ; /* IDA 0x19306; retn */
m->ip=0x9306;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_19307: ; /* IDA 0x19307; stc */
m->ip=0x9307;
su_skip_first_enter=0;
m->flags|=1;
goto L_19308;
L_19308: ; /* IDA 0x19308; retn */
m->ip=0x9308;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1EB6C: ; /* IDA 0x1EB6C; mov     ax, 12h */
m->ip=0xeb6c;
/* Original entry: sub_1EB6C */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0xeb6c,h->user);
m->ax=(uint16_t)(18);
goto L_1EB6F;
L_1EB6F: ; /* IDA 0x1EB6F; int     10h; - VIDEO - SET VIDEO MODE */
m->ip=0xeb6f;
su_skip_first_enter=0;
su_interrupt(m,16,0xeb71,h);
goto L_1EB71;
L_1EB71: ; /* IDA 0x1EB71; mov     dx, 3DAh */
m->ip=0xeb71;
su_skip_first_enter=0;
m->dx=(uint16_t)(986);
goto L_1EB74;
L_1EB74: ; /* IDA 0x1EB74; in      al, dx; Video status bits: */
m->ip=0xeb74;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(wolong_startup_in8(m->dx)));
goto L_1EB75;
L_1EB75: ; /* IDA 0x1EB75; mov     dx, 3C0h */
m->ip=0xeb75;
su_skip_first_enter=0;
m->dx=(uint16_t)(960);
goto L_1EB78;
L_1EB78: ; /* IDA 0x1EB78; mov     al, 14h */
m->ip=0xeb78;
su_skip_first_enter=0;
#if KI_STARTUP_MUTATION == 4
vg_al(m,(uint8_t)(21));
#else
vg_al(m,(uint8_t)(20));
#endif
goto L_1EB7A;
L_1EB7A: ; /* IDA 0x1EB7A; out     dx, al; EGA: unknown data port */
m->ip=0xeb7a;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EB7B;
L_1EB7B: ; /* IDA 0x1EB7B; xor     al, al */
m->ip=0xeb7b;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_1EB7D;
L_1EB7D: ; /* IDA 0x1EB7D; out     dx, al; EGA: palette register: select colors for attribute AL: */
m->ip=0xeb7d;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EB7E;
L_1EB7E: ; /* IDA 0x1EB7E; mov     cx, 10h */
m->ip=0xeb7e;
su_skip_first_enter=0;
m->cx=(uint16_t)(16);
goto L_1EB81;
L_1EB81: ; /* IDA 0x1EB81; xor     al, al */
m->ip=0xeb81;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_1EB83;
L_1EB83: ; /* IDA 0x1EB83; out     dx, al; EGA: palette register: select colors for attribute AL: */
m->ip=0xeb83;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EB84;
L_1EB84: ; /* IDA 0x1EB84; out     dx, al; EGA: palette register: select colors for attribute AL: */
m->ip=0xeb84;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EB85;
L_1EB85: ; /* IDA 0x1EB85; inc     al */
m->ip=0xeb85;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(ec_inc(m,(uint8_t)m->ax,8)));
goto L_1EB87;
L_1EB87: ; /* IDA 0x1EB87; loop    loc_1EB83 */
m->ip=0xeb87;
su_skip_first_enter=0;
if(--m->cx) {goto L_1EB83;}
goto L_1EB89;
L_1EB89: ; /* IDA 0x1EB89; mov     al, 11h */
m->ip=0xeb89;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(17));
goto L_1EB8B;
L_1EB8B: ; /* IDA 0x1EB8B; out     dx, al; EGA: select overscan (border) color. */
m->ip=0xeb8b;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EB8C;
L_1EB8C: ; /* IDA 0x1EB8C; xor     al, al */
m->ip=0xeb8c;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_1EB8E;
L_1EB8E: ; /* IDA 0x1EB8E; out     dx, al; EGA: palette register: select colors for attribute AL: */
m->ip=0xeb8e;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EB8F;
L_1EB8F: ; /* IDA 0x1EB8F; mov     al, 20h ; ' ' */
m->ip=0xeb8f;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(32));
goto L_1EB91;
L_1EB91: ; /* IDA 0x1EB91; out     dx, al; EGA: palette register: select colors for attribute AL: */
m->ip=0xeb91;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EB92;
L_1EB92: ; /* IDA 0x1EB92; retn */
m->ip=0xeb92;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
L_1EB93: ; /* IDA 0x1EB93; push    es */
m->ip=0xeb93;
/* Original entry: sub_1EB93 */
if(su_skip_first_enter)su_skip_first_enter=0;else if(h&&h->enter)h->enter(m,0xeb93,h->user);
ec_push(m,m->es);
goto L_1EB94;
L_1EB94: ; /* IDA 0x1EB94; push    cx */
m->ip=0xeb94;
su_skip_first_enter=0;
ec_push(m,m->cx);
goto L_1EB95;
L_1EB95: ; /* IDA 0x1EB95; push    di */
m->ip=0xeb95;
su_skip_first_enter=0;
ec_push(m,m->di);
goto L_1EB96;
L_1EB96: ; /* IDA 0x1EB96; mov     cx, 0A000h */
m->ip=0xeb96;
su_skip_first_enter=0;
m->cx=(uint16_t)(40960);
goto L_1EB99;
L_1EB99: ; /* IDA 0x1EB99; mov     es, cx */
m->ip=0xeb99;
su_skip_first_enter=0;
m->es=(uint16_t)(m->cx);
goto L_1EB9B;
L_1EB9B: ; /* IDA 0x1EB9B; mov     dx, 3DAh */
m->ip=0xeb9b;
su_skip_first_enter=0;
m->dx=(uint16_t)(986);
goto L_1EB9E;
L_1EB9E: ; /* IDA 0x1EB9E; in      al, dx; Video status bits: */
m->ip=0xeb9e;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(wolong_startup_in8(m->dx)));
goto L_1EB9F;
L_1EB9F: ; /* IDA 0x1EB9F; mov     dx, 3C0h */
m->ip=0xeb9f;
su_skip_first_enter=0;
m->dx=(uint16_t)(960);
goto L_1EBA2;
L_1EBA2: ; /* IDA 0x1EBA2; mov     al, 11h */
m->ip=0xeba2;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(17));
goto L_1EBA4;
L_1EBA4: ; /* IDA 0x1EBA4; out     dx, al; EGA: select overscan (border) color. */
m->ip=0xeba4;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EBA5;
L_1EBA5: ; /* IDA 0x1EBA5; mov     al, ah */
m->ip=0xeba5;
su_skip_first_enter=0;
vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
goto L_1EBA7;
L_1EBA7: ; /* IDA 0x1EBA7; out     dx, al; EGA: palette register: select colors for attribute AL: */
m->ip=0xeba7;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EBA8;
L_1EBA8: ; /* IDA 0x1EBA8; mov     al, 20h ; ' ' */
m->ip=0xeba8;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(32));
goto L_1EBAA;
L_1EBAA: ; /* IDA 0x1EBAA; out     dx, al; EGA: palette register: select colors for attribute AL: */
m->ip=0xebaa;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EBAB;
L_1EBAB: ; /* IDA 0x1EBAB; mov     dx, 3CEh */
m->ip=0xebab;
su_skip_first_enter=0;
m->dx=(uint16_t)(974);
goto L_1EBAE;
L_1EBAE: ; /* IDA 0x1EBAE; mov     al, 5 */
m->ip=0xebae;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(5));
goto L_1EBB0;
L_1EBB0: ; /* IDA 0x1EBB0; out     dx, al; EGA: graph 1 and 2 addr reg: */
m->ip=0xebb0;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EBB1;
L_1EBB1: ; /* IDA 0x1EBB1; inc     dx */
m->ip=0xebb1;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_inc(m,m->dx,16));
goto L_1EBB2;
L_1EBB2: ; /* IDA 0x1EBB2; mov     al, 3 */
m->ip=0xebb2;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(3));
goto L_1EBB4;
L_1EBB4: ; /* IDA 0x1EBB4; out     dx, al; EGA port: graphics controller data register */
m->ip=0xebb4;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EBB5;
L_1EBB5: ; /* IDA 0x1EBB5; dec     dx */
m->ip=0xebb5;
su_skip_first_enter=0;
m->dx=(uint16_t)(st_dec(m,m->dx,16));
goto L_1EBB6;
L_1EBB6: ; /* IDA 0x1EBB6; xor     al, al */
m->ip=0xebb6;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_1EBB8;
L_1EBB8: ; /* IDA 0x1EBB8; out     dx, al; EGA: graph 1 and 2 addr reg: */
m->ip=0xebb8;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EBB9;
L_1EBB9: ; /* IDA 0x1EBB9; inc     dx */
m->ip=0xebb9;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_inc(m,m->dx,16));
goto L_1EBBA;
L_1EBBA: ; /* IDA 0x1EBBA; mov     al, ah */
m->ip=0xebba;
su_skip_first_enter=0;
vg_al(m,(uint8_t)((uint8_t)(m->ax>>8)));
goto L_1EBBC;
L_1EBBC: ; /* IDA 0x1EBBC; out     dx, al; EGA port: graphics controller data register */
m->ip=0xebbc;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EBBD;
L_1EBBD: ; /* IDA 0x1EBBD; dec     dx */
m->ip=0xebbd;
su_skip_first_enter=0;
m->dx=(uint16_t)(st_dec(m,m->dx,16));
goto L_1EBBE;
L_1EBBE: ; /* IDA 0x1EBBE; mov     al, 3 */
m->ip=0xebbe;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(3));
goto L_1EBC0;
L_1EBC0: ; /* IDA 0x1EBC0; out     dx, al; EGA: graph 1 and 2 addr reg: */
m->ip=0xebc0;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EBC1;
L_1EBC1: ; /* IDA 0x1EBC1; inc     dx */
m->ip=0xebc1;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_inc(m,m->dx,16));
goto L_1EBC2;
L_1EBC2: ; /* IDA 0x1EBC2; xor     al, al */
m->ip=0xebc2;
su_skip_first_enter=0;
{uint16_t result=(((uint8_t)m->ax)^((uint8_t)m->ax))&255;ec_logic(m,result,8);vg_al(m,(uint8_t)(result));}
goto L_1EBC4;
L_1EBC4: ; /* IDA 0x1EBC4; out     dx, al; EGA port: graphics controller data register */
m->ip=0xebc4;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EBC5;
L_1EBC5: ; /* IDA 0x1EBC5; dec     dx */
m->ip=0xebc5;
su_skip_first_enter=0;
m->dx=(uint16_t)(st_dec(m,m->dx,16));
goto L_1EBC6;
L_1EBC6: ; /* IDA 0x1EBC6; mov     al, 8 */
m->ip=0xebc6;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(8));
goto L_1EBC8;
L_1EBC8: ; /* IDA 0x1EBC8; out     dx, al; EGA: graph 1 and 2 addr reg: */
m->ip=0xebc8;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EBC9;
L_1EBC9: ; /* IDA 0x1EBC9; inc     dx */
m->ip=0xebc9;
su_skip_first_enter=0;
m->dx=(uint16_t)(ec_inc(m,m->dx,16));
goto L_1EBCA;
L_1EBCA: ; /* IDA 0x1EBCA; mov     al, 0FFh */
m->ip=0xebca;
su_skip_first_enter=0;
vg_al(m,(uint8_t)(255));
goto L_1EBCC;
L_1EBCC: ; /* IDA 0x1EBCC; out     dx, al; EGA port: graphics controller data register */
m->ip=0xebcc;
su_skip_first_enter=0;
dl_out(m,m->dx,(uint8_t)m->ax,8);
goto L_1EBCD;
L_1EBCD: ; /* IDA 0x1EBCD; xor     di, di */
m->ip=0xebcd;
su_skip_first_enter=0;
{uint16_t result=((m->di)^(m->di))&65535;ec_logic(m,result,16);m->di=(uint16_t)(result);}
goto L_1EBCF;
L_1EBCF: ; /* IDA 0x1EBCF; mov     cx, 4B00h */
m->ip=0xebcf;
su_skip_first_enter=0;
m->cx=(uint16_t)(19200);
goto L_1EBD2;
L_1EBD2: ; /* IDA 0x1EBD2; cld */
m->ip=0xebd2;
su_skip_first_enter=0;
m->flags&=(uint16_t)~0x400u;
goto L_1EBD3;
L_1EBD3: ; /* IDA 0x1EBD3; mov     ax, 0FFFFh */
m->ip=0xebd3;
su_skip_first_enter=0;
m->ax=(uint16_t)(65535);
goto L_1EBD6;
L_1EBD6: ; /* IDA 0x1EBD6; rep stosw */
m->ip=0xebd6;
su_skip_first_enter=0;
mc_string(m,2,1);
goto L_1EBD8;
L_1EBD8: ; /* IDA 0x1EBD8; pop     di */
m->ip=0xebd8;
su_skip_first_enter=0;
m->di=(uint16_t)(ec_pop(m));
goto L_1EBD9;
L_1EBD9: ; /* IDA 0x1EBD9; pop     cx */
m->ip=0xebd9;
su_skip_first_enter=0;
m->cx=(uint16_t)(ec_pop(m));
goto L_1EBDA;
L_1EBDA: ; /* IDA 0x1EBDA; pop     es */
m->ip=0xebda;
su_skip_first_enter=0;
m->es=(uint16_t)(ec_pop(m));
goto L_1EBDB;
L_1EBDB: ; /* IDA 0x1EBDB; retn */
m->ip=0xebdb;
su_skip_first_enter=0;
m->ip=ec_pop(m);su_invoke_old=0;goto SU_DISPATCH;
}
int ki_startup_body(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {
if(!su_known(t))return 0;
uint16_t ss=m->ss,sp=m->sp,cs=m->cs,ret=dl_read16(m,ss,sp);
su_graph(m,t,h,1); /* Old invoke already emitted the entry hook. */
/* Old external ABI owns one final POP; reject abandoned native frames. */
if(t==0x356){uint16_t css=dl_read16(m,ss,(uint16_t)(sp+2));assert(m->ss==ss&&m->sp==(uint16_t)(sp+4)&&m->ip==ret&&m->cs==css);ec_push(m,css);}
else assert(m->ss==ss&&m->sp==(uint16_t)(sp+2)&&m->cs==cs&&m->ip==ret);
ec_push(m,ret);return 1;}
void ki_startup_invoke(KiMachine16 *m,uint16_t t,const KiEconomyHooks *h) {if(su_known(t))su_graph(m,t,h,0);else ki_army_invoke(m,t,h);}
void start(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x0,h,0);}
void sub_1005B(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x5b,h,0);}
void sub_1006B(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x6b,h,0);}
void sub_10155(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x155,h,0);}
void sub_10210(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x210,h,0);}
void sub_10236(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x236,h,0);}
void sub_1033B(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x33b,h,0);}
void sub_1034E(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x34e,h,0);}
void sub_10612(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x612,h,0);}
void sub_11A6E(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x1a6e,h,0);}
void sub_11AAE(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x1aae,h,0);}
void sub_11AC3(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x1ac3,h,0);}
void sub_11BE0(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x1be0,h,0);}
void sub_11CD0(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x1cd0,h,0);}
void sub_18B12(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x8b12,h,0);}
void sub_18B40(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x8b40,h,0);}
void sub_18CAE(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x8cae,h,0);}
void sub_18DC8(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x8dc8,h,0);}
void sub_18E5A(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x8e5a,h,0);}
void sub_18EA0(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x8ea0,h,0);}
void sub_18F37(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x8f37,h,0);}
void sub_18F6D(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x8f6d,h,0);}
void sub_18F7C(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x8f7c,h,0);}
void sub_18FC9(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x8fc9,h,0);}
void sub_1905E(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x905e,h,0);}
void sub_1906E(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x906e,h,0);}
void sub_1907F(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x907f,h,0);}
void sub_1908D(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x908d,h,0);}
void sub_190A3(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x90a3,h,0);}
void sub_190C0(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x90c0,h,0);}
void sub_1912D(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x912d,h,0);}
void sub_19144(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x9144,h,0);}
void sub_1915B(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x915b,h,0);}
void sub_1917C(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x917c,h,0);}
void sub_19214(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x9214,h,0);}
void sub_19223(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x9223,h,0);}
void sub_1928A(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x928a,h,0);}
void sub_192E3(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x92e3,h,0);}
void sub_192EE(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x92ee,h,0);}
void sub_192F5(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x92f5,h,0);}
void sub_1EB6C(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0xeb6c,h,0);}
void sub_1EB93(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0xeb93,h,0);}
void raw_entry_10043(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x43,h,0);}
void raw_entry_1006A(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x6a,h,0);}
void raw_entry_10356(KiMachine16 *m,const KiEconomyHooks *h) {su_graph(m,0x356,h,0);}

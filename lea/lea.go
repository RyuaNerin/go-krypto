package lea

func rol(W, i uint32) uint32 {
	return (((W) << (i)) | ((W) >> (32 - (i))))
}
func ror(W, i uint32) uint32 {
	return (((W) >> (i)) | ((W) << (32 - (i))))
}

func loadU32(v []byte, index int) uint32 {
	return (uint32(v[4*index+3]) << 24) | (uint32(v[4*index+2]) << 16) | (uint32(v[4*index+1]) << 8) | (uint32(v[4*index+0]))
}

func (key *lea_key) lea_set_key_generic(mk []byte, mk_len int) {
	switch mk_len {
	case 16:
		key.rk[0] = rol(loadU32(mk, 0)+delta[0][0], 1)
		key.rk[6] = rol(key.rk[0]+delta[1][1], 1)
		key.rk[12] = rol(key.rk[6]+delta[2][2], 1)
		key.rk[18] = rol(key.rk[12]+delta[3][3], 1)
		key.rk[24] = rol(key.rk[18]+delta[0][4], 1)
		key.rk[30] = rol(key.rk[24]+delta[1][5], 1)
		key.rk[36] = rol(key.rk[30]+delta[2][6], 1)
		key.rk[42] = rol(key.rk[36]+delta[3][7], 1)
		key.rk[48] = rol(key.rk[42]+delta[0][8], 1)
		key.rk[54] = rol(key.rk[48]+delta[1][9], 1)
		key.rk[60] = rol(key.rk[54]+delta[2][10], 1)
		key.rk[66] = rol(key.rk[60]+delta[3][11], 1)
		key.rk[72] = rol(key.rk[66]+delta[0][12], 1)
		key.rk[78] = rol(key.rk[72]+delta[1][13], 1)
		key.rk[84] = rol(key.rk[78]+delta[2][14], 1)
		key.rk[90] = rol(key.rk[84]+delta[3][15], 1)
		key.rk[96] = rol(key.rk[90]+delta[0][16], 1)
		key.rk[102] = rol(key.rk[96]+delta[1][17], 1)
		key.rk[108] = rol(key.rk[102]+delta[2][18], 1)
		key.rk[114] = rol(key.rk[108]+delta[3][19], 1)
		key.rk[120] = rol(key.rk[114]+delta[0][20], 1)
		key.rk[126] = rol(key.rk[120]+delta[1][21], 1)
		key.rk[132] = rol(key.rk[126]+delta[2][22], 1)
		key.rk[138] = rol(key.rk[132]+delta[3][23], 1)

		/**
		key.rk[  1] = key.rk[  3] = key.rk[  5] = ROL(loadU32(mk,1) + delta[0][ 1], 3);
		key.rk[  7] = key.rk[  9] = key.rk[ 11] = ROL(key.rk[  1] + delta[1][ 2], 3);
		key.rk[ 13] = key.rk[ 15] = key.rk[ 17] = ROL(key.rk[  7] + delta[2][ 3], 3);
		key.rk[ 19] = key.rk[ 21] = key.rk[ 23] = ROL(key.rk[ 13] + delta[3][ 4], 3);
		key.rk[ 25] = key.rk[ 27] = key.rk[ 29] = ROL(key.rk[ 19] + delta[0][ 5], 3);
		key.rk[ 31] = key.rk[ 33] = key.rk[ 35] = ROL(key.rk[ 25] + delta[1][ 6], 3);
		key.rk[ 37] = key.rk[ 39] = key.rk[ 41] = ROL(key.rk[ 31] + delta[2][ 7], 3);
		key.rk[ 43] = key.rk[ 45] = key.rk[ 47] = ROL(key.rk[ 37] + delta[3][ 8], 3);
		key.rk[ 49] = key.rk[ 51] = key.rk[ 53] = ROL(key.rk[ 43] + delta[0][ 9], 3);
		key.rk[ 55] = key.rk[ 57] = key.rk[ 59] = ROL(key.rk[ 49] + delta[1][10], 3);
		key.rk[ 61] = key.rk[ 63] = key.rk[ 65] = ROL(key.rk[ 55] + delta[2][11], 3);
		key.rk[ 67] = key.rk[ 69] = key.rk[ 71] = ROL(key.rk[ 61] + delta[3][12], 3);
		key.rk[ 73] = key.rk[ 75] = key.rk[ 77] = ROL(key.rk[ 67] + delta[0][13], 3);
		key.rk[ 79] = key.rk[ 81] = key.rk[ 83] = ROL(key.rk[ 73] + delta[1][14], 3);
		key.rk[ 85] = key.rk[ 87] = key.rk[ 89] = ROL(key.rk[ 79] + delta[2][15], 3);
		key.rk[ 91] = key.rk[ 93] = key.rk[ 95] = ROL(key.rk[ 85] + delta[3][16], 3);
		key.rk[ 97] = key.rk[ 99] = key.rk[101] = ROL(key.rk[ 91] + delta[0][17], 3);
		key.rk[103] = key.rk[105] = key.rk[107] = ROL(key.rk[ 97] + delta[1][18], 3);
		key.rk[109] = key.rk[111] = key.rk[113] = ROL(key.rk[103] + delta[2][19], 3);
		key.rk[115] = key.rk[117] = key.rk[119] = ROL(key.rk[109] + delta[3][20], 3);
		key.rk[121] = key.rk[123] = key.rk[125] = ROL(key.rk[115] + delta[0][21], 3);
		key.rk[127] = key.rk[129] = key.rk[131] = ROL(key.rk[121] + delta[1][22], 3);
		key.rk[133] = key.rk[135] = key.rk[137] = ROL(key.rk[127] + delta[2][23], 3);
		key.rk[139] = key.rk[141] = key.rk[143] = ROL(key.rk[133] + delta[3][24], 3);
		*/
		var tmp uint32
		tmp = rol(loadU32(mk, 1)+delta[0][1], 3)
		key.rk[1] = tmp
		key.rk[3] = tmp
		key.rk[5] = tmp

		for i := 1; i <= 23; i++ {
			tmp = rol(key.rk[(i-1)*6+1]+delta[i%4][i+1], 3)

			key.rk[i*6+1] = tmp
			key.rk[i*6+3] = tmp
			key.rk[i*6+5] = tmp
		}

		key.rk[2] = rol(loadU32(mk, 2)+delta[0][2], 6)
		key.rk[8] = rol(key.rk[2]+delta[1][3], 6)
		key.rk[14] = rol(key.rk[8]+delta[2][4], 6)
		key.rk[20] = rol(key.rk[14]+delta[3][5], 6)
		key.rk[26] = rol(key.rk[20]+delta[0][6], 6)
		key.rk[32] = rol(key.rk[26]+delta[1][7], 6)
		key.rk[38] = rol(key.rk[32]+delta[2][8], 6)
		key.rk[44] = rol(key.rk[38]+delta[3][9], 6)
		key.rk[50] = rol(key.rk[44]+delta[0][10], 6)
		key.rk[56] = rol(key.rk[50]+delta[1][11], 6)
		key.rk[62] = rol(key.rk[56]+delta[2][12], 6)
		key.rk[68] = rol(key.rk[62]+delta[3][13], 6)
		key.rk[74] = rol(key.rk[68]+delta[0][14], 6)
		key.rk[80] = rol(key.rk[74]+delta[1][15], 6)
		key.rk[86] = rol(key.rk[80]+delta[2][16], 6)
		key.rk[92] = rol(key.rk[86]+delta[3][17], 6)
		key.rk[98] = rol(key.rk[92]+delta[0][18], 6)
		key.rk[104] = rol(key.rk[98]+delta[1][19], 6)
		key.rk[110] = rol(key.rk[104]+delta[2][20], 6)
		key.rk[116] = rol(key.rk[110]+delta[3][21], 6)
		key.rk[122] = rol(key.rk[116]+delta[0][22], 6)
		key.rk[128] = rol(key.rk[122]+delta[1][23], 6)
		key.rk[134] = rol(key.rk[128]+delta[2][24], 6)
		key.rk[140] = rol(key.rk[134]+delta[3][25], 6)

		key.rk[4] = rol(loadU32(mk, 3)+delta[0][3], 11)
		key.rk[10] = rol(key.rk[4]+delta[1][4], 11)
		key.rk[16] = rol(key.rk[10]+delta[2][5], 11)
		key.rk[22] = rol(key.rk[16]+delta[3][6], 11)
		key.rk[28] = rol(key.rk[22]+delta[0][7], 11)
		key.rk[34] = rol(key.rk[28]+delta[1][8], 11)
		key.rk[40] = rol(key.rk[34]+delta[2][9], 11)
		key.rk[46] = rol(key.rk[40]+delta[3][10], 11)
		key.rk[52] = rol(key.rk[46]+delta[0][11], 11)
		key.rk[58] = rol(key.rk[52]+delta[1][12], 11)
		key.rk[64] = rol(key.rk[58]+delta[2][13], 11)
		key.rk[70] = rol(key.rk[64]+delta[3][14], 11)
		key.rk[76] = rol(key.rk[70]+delta[0][15], 11)
		key.rk[82] = rol(key.rk[76]+delta[1][16], 11)
		key.rk[88] = rol(key.rk[82]+delta[2][17], 11)
		key.rk[94] = rol(key.rk[88]+delta[3][18], 11)
		key.rk[100] = rol(key.rk[94]+delta[0][19], 11)
		key.rk[106] = rol(key.rk[100]+delta[1][20], 11)
		key.rk[112] = rol(key.rk[106]+delta[2][21], 11)
		key.rk[118] = rol(key.rk[112]+delta[3][22], 11)
		key.rk[124] = rol(key.rk[118]+delta[0][23], 11)
		key.rk[130] = rol(key.rk[124]+delta[1][24], 11)
		key.rk[136] = rol(key.rk[130]+delta[2][25], 11)
		key.rk[142] = rol(key.rk[136]+delta[3][26], 11)
		break

	case 24:
		key.rk[0] = rol(loadU32(mk, 0)+delta[0][0], 1)
		key.rk[6] = rol(key.rk[0]+delta[1][1], 1)
		key.rk[12] = rol(key.rk[6]+delta[2][2], 1)
		key.rk[18] = rol(key.rk[12]+delta[3][3], 1)
		key.rk[24] = rol(key.rk[18]+delta[4][4], 1)
		key.rk[30] = rol(key.rk[24]+delta[5][5], 1)
		key.rk[36] = rol(key.rk[30]+delta[0][6], 1)
		key.rk[42] = rol(key.rk[36]+delta[1][7], 1)
		key.rk[48] = rol(key.rk[42]+delta[2][8], 1)
		key.rk[54] = rol(key.rk[48]+delta[3][9], 1)
		key.rk[60] = rol(key.rk[54]+delta[4][10], 1)
		key.rk[66] = rol(key.rk[60]+delta[5][11], 1)
		key.rk[72] = rol(key.rk[66]+delta[0][12], 1)
		key.rk[78] = rol(key.rk[72]+delta[1][13], 1)
		key.rk[84] = rol(key.rk[78]+delta[2][14], 1)
		key.rk[90] = rol(key.rk[84]+delta[3][15], 1)
		key.rk[96] = rol(key.rk[90]+delta[4][16], 1)
		key.rk[102] = rol(key.rk[96]+delta[5][17], 1)
		key.rk[108] = rol(key.rk[102]+delta[0][18], 1)
		key.rk[114] = rol(key.rk[108]+delta[1][19], 1)
		key.rk[120] = rol(key.rk[114]+delta[2][20], 1)
		key.rk[126] = rol(key.rk[120]+delta[3][21], 1)
		key.rk[132] = rol(key.rk[126]+delta[4][22], 1)
		key.rk[138] = rol(key.rk[132]+delta[5][23], 1)
		key.rk[144] = rol(key.rk[138]+delta[0][24], 1)
		key.rk[150] = rol(key.rk[144]+delta[1][25], 1)
		key.rk[156] = rol(key.rk[150]+delta[2][26], 1)
		key.rk[162] = rol(key.rk[156]+delta[3][27], 1)

		key.rk[1] = rol(loadU32(mk, 1)+delta[0][1], 3)
		key.rk[7] = rol(key.rk[1]+delta[1][2], 3)
		key.rk[13] = rol(key.rk[7]+delta[2][3], 3)
		key.rk[19] = rol(key.rk[13]+delta[3][4], 3)
		key.rk[25] = rol(key.rk[19]+delta[4][5], 3)
		key.rk[31] = rol(key.rk[25]+delta[5][6], 3)
		key.rk[37] = rol(key.rk[31]+delta[0][7], 3)
		key.rk[43] = rol(key.rk[37]+delta[1][8], 3)
		key.rk[49] = rol(key.rk[43]+delta[2][9], 3)
		key.rk[55] = rol(key.rk[49]+delta[3][10], 3)
		key.rk[61] = rol(key.rk[55]+delta[4][11], 3)
		key.rk[67] = rol(key.rk[61]+delta[5][12], 3)
		key.rk[73] = rol(key.rk[67]+delta[0][13], 3)
		key.rk[79] = rol(key.rk[73]+delta[1][14], 3)
		key.rk[85] = rol(key.rk[79]+delta[2][15], 3)
		key.rk[91] = rol(key.rk[85]+delta[3][16], 3)
		key.rk[97] = rol(key.rk[91]+delta[4][17], 3)
		key.rk[103] = rol(key.rk[97]+delta[5][18], 3)
		key.rk[109] = rol(key.rk[103]+delta[0][19], 3)
		key.rk[115] = rol(key.rk[109]+delta[1][20], 3)
		key.rk[121] = rol(key.rk[115]+delta[2][21], 3)
		key.rk[127] = rol(key.rk[121]+delta[3][22], 3)
		key.rk[133] = rol(key.rk[127]+delta[4][23], 3)
		key.rk[139] = rol(key.rk[133]+delta[5][24], 3)
		key.rk[145] = rol(key.rk[139]+delta[0][25], 3)
		key.rk[151] = rol(key.rk[145]+delta[1][26], 3)
		key.rk[157] = rol(key.rk[151]+delta[2][27], 3)
		key.rk[163] = rol(key.rk[157]+delta[3][28], 3)

		key.rk[2] = rol(loadU32(mk, 2)+delta[0][2], 6)
		key.rk[8] = rol(key.rk[2]+delta[1][3], 6)
		key.rk[14] = rol(key.rk[8]+delta[2][4], 6)
		key.rk[20] = rol(key.rk[14]+delta[3][5], 6)
		key.rk[26] = rol(key.rk[20]+delta[4][6], 6)
		key.rk[32] = rol(key.rk[26]+delta[5][7], 6)
		key.rk[38] = rol(key.rk[32]+delta[0][8], 6)
		key.rk[44] = rol(key.rk[38]+delta[1][9], 6)
		key.rk[50] = rol(key.rk[44]+delta[2][10], 6)
		key.rk[56] = rol(key.rk[50]+delta[3][11], 6)
		key.rk[62] = rol(key.rk[56]+delta[4][12], 6)
		key.rk[68] = rol(key.rk[62]+delta[5][13], 6)
		key.rk[74] = rol(key.rk[68]+delta[0][14], 6)
		key.rk[80] = rol(key.rk[74]+delta[1][15], 6)
		key.rk[86] = rol(key.rk[80]+delta[2][16], 6)
		key.rk[92] = rol(key.rk[86]+delta[3][17], 6)
		key.rk[98] = rol(key.rk[92]+delta[4][18], 6)
		key.rk[104] = rol(key.rk[98]+delta[5][19], 6)
		key.rk[110] = rol(key.rk[104]+delta[0][20], 6)
		key.rk[116] = rol(key.rk[110]+delta[1][21], 6)
		key.rk[122] = rol(key.rk[116]+delta[2][22], 6)
		key.rk[128] = rol(key.rk[122]+delta[3][23], 6)
		key.rk[134] = rol(key.rk[128]+delta[4][24], 6)
		key.rk[140] = rol(key.rk[134]+delta[5][25], 6)
		key.rk[146] = rol(key.rk[140]+delta[0][26], 6)
		key.rk[152] = rol(key.rk[146]+delta[1][27], 6)
		key.rk[158] = rol(key.rk[152]+delta[2][28], 6)
		key.rk[164] = rol(key.rk[158]+delta[3][29], 6)

		key.rk[3] = rol(loadU32(mk, 3)+delta[0][3], 11)
		key.rk[9] = rol(key.rk[3]+delta[1][4], 11)
		key.rk[15] = rol(key.rk[9]+delta[2][5], 11)
		key.rk[21] = rol(key.rk[15]+delta[3][6], 11)
		key.rk[27] = rol(key.rk[21]+delta[4][7], 11)
		key.rk[33] = rol(key.rk[27]+delta[5][8], 11)
		key.rk[39] = rol(key.rk[33]+delta[0][9], 11)
		key.rk[45] = rol(key.rk[39]+delta[1][10], 11)
		key.rk[51] = rol(key.rk[45]+delta[2][11], 11)
		key.rk[57] = rol(key.rk[51]+delta[3][12], 11)
		key.rk[63] = rol(key.rk[57]+delta[4][13], 11)
		key.rk[69] = rol(key.rk[63]+delta[5][14], 11)
		key.rk[75] = rol(key.rk[69]+delta[0][15], 11)
		key.rk[81] = rol(key.rk[75]+delta[1][16], 11)
		key.rk[87] = rol(key.rk[81]+delta[2][17], 11)
		key.rk[93] = rol(key.rk[87]+delta[3][18], 11)
		key.rk[99] = rol(key.rk[93]+delta[4][19], 11)
		key.rk[105] = rol(key.rk[99]+delta[5][20], 11)
		key.rk[111] = rol(key.rk[105]+delta[0][21], 11)
		key.rk[117] = rol(key.rk[111]+delta[1][22], 11)
		key.rk[123] = rol(key.rk[117]+delta[2][23], 11)
		key.rk[129] = rol(key.rk[123]+delta[3][24], 11)
		key.rk[135] = rol(key.rk[129]+delta[4][25], 11)
		key.rk[141] = rol(key.rk[135]+delta[5][26], 11)
		key.rk[147] = rol(key.rk[141]+delta[0][27], 11)
		key.rk[153] = rol(key.rk[147]+delta[1][28], 11)
		key.rk[159] = rol(key.rk[153]+delta[2][29], 11)
		key.rk[165] = rol(key.rk[159]+delta[3][30], 11)

		key.rk[4] = rol(loadU32(mk, 4)+delta[0][4], 13)
		key.rk[10] = rol(key.rk[4]+delta[1][5], 13)
		key.rk[16] = rol(key.rk[10]+delta[2][6], 13)
		key.rk[22] = rol(key.rk[16]+delta[3][7], 13)
		key.rk[28] = rol(key.rk[22]+delta[4][8], 13)
		key.rk[34] = rol(key.rk[28]+delta[5][9], 13)
		key.rk[40] = rol(key.rk[34]+delta[0][10], 13)
		key.rk[46] = rol(key.rk[40]+delta[1][11], 13)
		key.rk[52] = rol(key.rk[46]+delta[2][12], 13)
		key.rk[58] = rol(key.rk[52]+delta[3][13], 13)
		key.rk[64] = rol(key.rk[58]+delta[4][14], 13)
		key.rk[70] = rol(key.rk[64]+delta[5][15], 13)
		key.rk[76] = rol(key.rk[70]+delta[0][16], 13)
		key.rk[82] = rol(key.rk[76]+delta[1][17], 13)
		key.rk[88] = rol(key.rk[82]+delta[2][18], 13)
		key.rk[94] = rol(key.rk[88]+delta[3][19], 13)
		key.rk[100] = rol(key.rk[94]+delta[4][20], 13)
		key.rk[106] = rol(key.rk[100]+delta[5][21], 13)
		key.rk[112] = rol(key.rk[106]+delta[0][22], 13)
		key.rk[118] = rol(key.rk[112]+delta[1][23], 13)
		key.rk[124] = rol(key.rk[118]+delta[2][24], 13)
		key.rk[130] = rol(key.rk[124]+delta[3][25], 13)
		key.rk[136] = rol(key.rk[130]+delta[4][26], 13)
		key.rk[142] = rol(key.rk[136]+delta[5][27], 13)
		key.rk[148] = rol(key.rk[142]+delta[0][28], 13)
		key.rk[154] = rol(key.rk[148]+delta[1][29], 13)
		key.rk[160] = rol(key.rk[154]+delta[2][30], 13)
		key.rk[166] = rol(key.rk[160]+delta[3][31], 13)

		key.rk[5] = rol(loadU32(mk, 5)+delta[0][5], 17)
		key.rk[11] = rol(key.rk[5]+delta[1][6], 17)
		key.rk[17] = rol(key.rk[11]+delta[2][7], 17)
		key.rk[23] = rol(key.rk[17]+delta[3][8], 17)
		key.rk[29] = rol(key.rk[23]+delta[4][9], 17)
		key.rk[35] = rol(key.rk[29]+delta[5][10], 17)
		key.rk[41] = rol(key.rk[35]+delta[0][11], 17)
		key.rk[47] = rol(key.rk[41]+delta[1][12], 17)
		key.rk[53] = rol(key.rk[47]+delta[2][13], 17)
		key.rk[59] = rol(key.rk[53]+delta[3][14], 17)
		key.rk[65] = rol(key.rk[59]+delta[4][15], 17)
		key.rk[71] = rol(key.rk[65]+delta[5][16], 17)
		key.rk[77] = rol(key.rk[71]+delta[0][17], 17)
		key.rk[83] = rol(key.rk[77]+delta[1][18], 17)
		key.rk[89] = rol(key.rk[83]+delta[2][19], 17)
		key.rk[95] = rol(key.rk[89]+delta[3][20], 17)
		key.rk[101] = rol(key.rk[95]+delta[4][21], 17)
		key.rk[107] = rol(key.rk[101]+delta[5][22], 17)
		key.rk[113] = rol(key.rk[107]+delta[0][23], 17)
		key.rk[119] = rol(key.rk[113]+delta[1][24], 17)
		key.rk[125] = rol(key.rk[119]+delta[2][25], 17)
		key.rk[131] = rol(key.rk[125]+delta[3][26], 17)
		key.rk[137] = rol(key.rk[131]+delta[4][27], 17)
		key.rk[143] = rol(key.rk[137]+delta[5][28], 17)
		key.rk[149] = rol(key.rk[143]+delta[0][29], 17)
		key.rk[155] = rol(key.rk[149]+delta[1][30], 17)
		key.rk[161] = rol(key.rk[155]+delta[2][31], 17)
		key.rk[167] = rol(key.rk[161]+delta[3][0], 17)
		break

	case 32:
		key.rk[0] = rol(loadU32(mk, 0)+delta[0][0], 1)
		key.rk[8] = rol(key.rk[0]+delta[1][3], 6)
		key.rk[16] = rol(key.rk[8]+delta[2][6], 13)
		key.rk[24] = rol(key.rk[16]+delta[4][4], 1)
		key.rk[32] = rol(key.rk[24]+delta[5][7], 6)
		key.rk[40] = rol(key.rk[32]+delta[6][10], 13)
		key.rk[48] = rol(key.rk[40]+delta[0][8], 1)
		key.rk[56] = rol(key.rk[48]+delta[1][11], 6)
		key.rk[64] = rol(key.rk[56]+delta[2][14], 13)
		key.rk[72] = rol(key.rk[64]+delta[4][12], 1)
		key.rk[80] = rol(key.rk[72]+delta[5][15], 6)
		key.rk[88] = rol(key.rk[80]+delta[6][18], 13)
		key.rk[96] = rol(key.rk[88]+delta[0][16], 1)
		key.rk[104] = rol(key.rk[96]+delta[1][19], 6)
		key.rk[112] = rol(key.rk[104]+delta[2][22], 13)
		key.rk[120] = rol(key.rk[112]+delta[4][20], 1)
		key.rk[128] = rol(key.rk[120]+delta[5][23], 6)
		key.rk[136] = rol(key.rk[128]+delta[6][26], 13)
		key.rk[144] = rol(key.rk[136]+delta[0][24], 1)
		key.rk[152] = rol(key.rk[144]+delta[1][27], 6)
		key.rk[160] = rol(key.rk[152]+delta[2][30], 13)
		key.rk[168] = rol(key.rk[160]+delta[4][28], 1)
		key.rk[176] = rol(key.rk[168]+delta[5][31], 6)
		key.rk[184] = rol(key.rk[176]+delta[6][2], 13)

		key.rk[1] = rol(loadU32(mk, 1)+delta[0][1], 3)
		key.rk[9] = rol(key.rk[1]+delta[1][4], 11)
		key.rk[17] = rol(key.rk[9]+delta[2][7], 17)
		key.rk[25] = rol(key.rk[17]+delta[4][5], 3)
		key.rk[33] = rol(key.rk[25]+delta[5][8], 11)
		key.rk[41] = rol(key.rk[33]+delta[6][11], 17)
		key.rk[49] = rol(key.rk[41]+delta[0][9], 3)
		key.rk[57] = rol(key.rk[49]+delta[1][12], 11)
		key.rk[65] = rol(key.rk[57]+delta[2][15], 17)
		key.rk[73] = rol(key.rk[65]+delta[4][13], 3)
		key.rk[81] = rol(key.rk[73]+delta[5][16], 11)
		key.rk[89] = rol(key.rk[81]+delta[6][19], 17)
		key.rk[97] = rol(key.rk[89]+delta[0][17], 3)
		key.rk[105] = rol(key.rk[97]+delta[1][20], 11)
		key.rk[113] = rol(key.rk[105]+delta[2][23], 17)
		key.rk[121] = rol(key.rk[113]+delta[4][21], 3)
		key.rk[129] = rol(key.rk[121]+delta[5][24], 11)
		key.rk[137] = rol(key.rk[129]+delta[6][27], 17)
		key.rk[145] = rol(key.rk[137]+delta[0][25], 3)
		key.rk[153] = rol(key.rk[145]+delta[1][28], 11)
		key.rk[161] = rol(key.rk[153]+delta[2][31], 17)
		key.rk[169] = rol(key.rk[161]+delta[4][29], 3)
		key.rk[177] = rol(key.rk[169]+delta[5][0], 11)
		key.rk[185] = rol(key.rk[177]+delta[6][3], 17)

		key.rk[2] = rol(loadU32(mk, 2)+delta[0][2], 6)
		key.rk[10] = rol(key.rk[2]+delta[1][5], 13)
		key.rk[18] = rol(key.rk[10]+delta[3][3], 1)
		key.rk[26] = rol(key.rk[18]+delta[4][6], 6)
		key.rk[34] = rol(key.rk[26]+delta[5][9], 13)
		key.rk[42] = rol(key.rk[34]+delta[7][7], 1)
		key.rk[50] = rol(key.rk[42]+delta[0][10], 6)
		key.rk[58] = rol(key.rk[50]+delta[1][13], 13)
		key.rk[66] = rol(key.rk[58]+delta[3][11], 1)
		key.rk[74] = rol(key.rk[66]+delta[4][14], 6)
		key.rk[82] = rol(key.rk[74]+delta[5][17], 13)
		key.rk[90] = rol(key.rk[82]+delta[7][15], 1)
		key.rk[98] = rol(key.rk[90]+delta[0][18], 6)
		key.rk[106] = rol(key.rk[98]+delta[1][21], 13)
		key.rk[114] = rol(key.rk[106]+delta[3][19], 1)
		key.rk[122] = rol(key.rk[114]+delta[4][22], 6)
		key.rk[130] = rol(key.rk[122]+delta[5][25], 13)
		key.rk[138] = rol(key.rk[130]+delta[7][23], 1)
		key.rk[146] = rol(key.rk[138]+delta[0][26], 6)
		key.rk[154] = rol(key.rk[146]+delta[1][29], 13)
		key.rk[162] = rol(key.rk[154]+delta[3][27], 1)
		key.rk[170] = rol(key.rk[162]+delta[4][30], 6)
		key.rk[178] = rol(key.rk[170]+delta[5][1], 13)
		key.rk[186] = rol(key.rk[178]+delta[7][31], 1)

		key.rk[3] = rol(loadU32(mk, 3)+delta[0][3], 11)
		key.rk[11] = rol(key.rk[3]+delta[1][6], 17)
		key.rk[19] = rol(key.rk[11]+delta[3][4], 3)
		key.rk[27] = rol(key.rk[19]+delta[4][7], 11)
		key.rk[35] = rol(key.rk[27]+delta[5][10], 17)
		key.rk[43] = rol(key.rk[35]+delta[7][8], 3)
		key.rk[51] = rol(key.rk[43]+delta[0][11], 11)
		key.rk[59] = rol(key.rk[51]+delta[1][14], 17)
		key.rk[67] = rol(key.rk[59]+delta[3][12], 3)
		key.rk[75] = rol(key.rk[67]+delta[4][15], 11)
		key.rk[83] = rol(key.rk[75]+delta[5][18], 17)
		key.rk[91] = rol(key.rk[83]+delta[7][16], 3)
		key.rk[99] = rol(key.rk[91]+delta[0][19], 11)
		key.rk[107] = rol(key.rk[99]+delta[1][22], 17)
		key.rk[115] = rol(key.rk[107]+delta[3][20], 3)
		key.rk[123] = rol(key.rk[115]+delta[4][23], 11)
		key.rk[131] = rol(key.rk[123]+delta[5][26], 17)
		key.rk[139] = rol(key.rk[131]+delta[7][24], 3)
		key.rk[147] = rol(key.rk[139]+delta[0][27], 11)
		key.rk[155] = rol(key.rk[147]+delta[1][30], 17)
		key.rk[163] = rol(key.rk[155]+delta[3][28], 3)
		key.rk[171] = rol(key.rk[163]+delta[4][31], 11)
		key.rk[179] = rol(key.rk[171]+delta[5][2], 17)
		key.rk[187] = rol(key.rk[179]+delta[7][0], 3)

		key.rk[4] = rol(loadU32(mk, 4)+delta[0][4], 13)
		key.rk[12] = rol(key.rk[4]+delta[2][2], 1)
		key.rk[20] = rol(key.rk[12]+delta[3][5], 6)
		key.rk[28] = rol(key.rk[20]+delta[4][8], 13)
		key.rk[36] = rol(key.rk[28]+delta[6][6], 1)
		key.rk[44] = rol(key.rk[36]+delta[7][9], 6)
		key.rk[52] = rol(key.rk[44]+delta[0][12], 13)
		key.rk[60] = rol(key.rk[52]+delta[2][10], 1)
		key.rk[68] = rol(key.rk[60]+delta[3][13], 6)
		key.rk[76] = rol(key.rk[68]+delta[4][16], 13)
		key.rk[84] = rol(key.rk[76]+delta[6][14], 1)
		key.rk[92] = rol(key.rk[84]+delta[7][17], 6)
		key.rk[100] = rol(key.rk[92]+delta[0][20], 13)
		key.rk[108] = rol(key.rk[100]+delta[2][18], 1)
		key.rk[116] = rol(key.rk[108]+delta[3][21], 6)
		key.rk[124] = rol(key.rk[116]+delta[4][24], 13)
		key.rk[132] = rol(key.rk[124]+delta[6][22], 1)
		key.rk[140] = rol(key.rk[132]+delta[7][25], 6)
		key.rk[148] = rol(key.rk[140]+delta[0][28], 13)
		key.rk[156] = rol(key.rk[148]+delta[2][26], 1)
		key.rk[164] = rol(key.rk[156]+delta[3][29], 6)
		key.rk[172] = rol(key.rk[164]+delta[4][0], 13)
		key.rk[180] = rol(key.rk[172]+delta[6][30], 1)
		key.rk[188] = rol(key.rk[180]+delta[7][1], 6)

		key.rk[5] = rol(loadU32(mk, 5)+delta[0][5], 17)
		key.rk[13] = rol(key.rk[5]+delta[2][3], 3)
		key.rk[21] = rol(key.rk[13]+delta[3][6], 11)
		key.rk[29] = rol(key.rk[21]+delta[4][9], 17)
		key.rk[37] = rol(key.rk[29]+delta[6][7], 3)
		key.rk[45] = rol(key.rk[37]+delta[7][10], 11)
		key.rk[53] = rol(key.rk[45]+delta[0][13], 17)
		key.rk[61] = rol(key.rk[53]+delta[2][11], 3)
		key.rk[69] = rol(key.rk[61]+delta[3][14], 11)
		key.rk[77] = rol(key.rk[69]+delta[4][17], 17)
		key.rk[85] = rol(key.rk[77]+delta[6][15], 3)
		key.rk[93] = rol(key.rk[85]+delta[7][18], 11)
		key.rk[101] = rol(key.rk[93]+delta[0][21], 17)
		key.rk[109] = rol(key.rk[101]+delta[2][19], 3)
		key.rk[117] = rol(key.rk[109]+delta[3][22], 11)
		key.rk[125] = rol(key.rk[117]+delta[4][25], 17)
		key.rk[133] = rol(key.rk[125]+delta[6][23], 3)
		key.rk[141] = rol(key.rk[133]+delta[7][26], 11)
		key.rk[149] = rol(key.rk[141]+delta[0][29], 17)
		key.rk[157] = rol(key.rk[149]+delta[2][27], 3)
		key.rk[165] = rol(key.rk[157]+delta[3][30], 11)
		key.rk[173] = rol(key.rk[165]+delta[4][1], 17)
		key.rk[181] = rol(key.rk[173]+delta[6][31], 3)
		key.rk[189] = rol(key.rk[181]+delta[7][2], 11)

		key.rk[6] = rol(loadU32(mk, 6)+delta[1][1], 1)
		key.rk[14] = rol(key.rk[6]+delta[2][4], 6)
		key.rk[22] = rol(key.rk[14]+delta[3][7], 13)
		key.rk[30] = rol(key.rk[22]+delta[5][5], 1)
		key.rk[38] = rol(key.rk[30]+delta[6][8], 6)
		key.rk[46] = rol(key.rk[38]+delta[7][11], 13)
		key.rk[54] = rol(key.rk[46]+delta[1][9], 1)
		key.rk[62] = rol(key.rk[54]+delta[2][12], 6)
		key.rk[70] = rol(key.rk[62]+delta[3][15], 13)
		key.rk[78] = rol(key.rk[70]+delta[5][13], 1)
		key.rk[86] = rol(key.rk[78]+delta[6][16], 6)
		key.rk[94] = rol(key.rk[86]+delta[7][19], 13)
		key.rk[102] = rol(key.rk[94]+delta[1][17], 1)
		key.rk[110] = rol(key.rk[102]+delta[2][20], 6)
		key.rk[118] = rol(key.rk[110]+delta[3][23], 13)
		key.rk[126] = rol(key.rk[118]+delta[5][21], 1)
		key.rk[134] = rol(key.rk[126]+delta[6][24], 6)
		key.rk[142] = rol(key.rk[134]+delta[7][27], 13)
		key.rk[150] = rol(key.rk[142]+delta[1][25], 1)
		key.rk[158] = rol(key.rk[150]+delta[2][28], 6)
		key.rk[166] = rol(key.rk[158]+delta[3][31], 13)
		key.rk[174] = rol(key.rk[166]+delta[5][29], 1)
		key.rk[182] = rol(key.rk[174]+delta[6][0], 6)
		key.rk[190] = rol(key.rk[182]+delta[7][3], 13)

		key.rk[7] = rol(loadU32(mk, 7)+delta[1][2], 3)
		key.rk[15] = rol(key.rk[7]+delta[2][5], 11)
		key.rk[23] = rol(key.rk[15]+delta[3][8], 17)
		key.rk[31] = rol(key.rk[23]+delta[5][6], 3)
		key.rk[39] = rol(key.rk[31]+delta[6][9], 11)
		key.rk[47] = rol(key.rk[39]+delta[7][12], 17)
		key.rk[55] = rol(key.rk[47]+delta[1][10], 3)
		key.rk[63] = rol(key.rk[55]+delta[2][13], 11)
		key.rk[71] = rol(key.rk[63]+delta[3][16], 17)
		key.rk[79] = rol(key.rk[71]+delta[5][14], 3)
		key.rk[87] = rol(key.rk[79]+delta[6][17], 11)
		key.rk[95] = rol(key.rk[87]+delta[7][20], 17)
		key.rk[103] = rol(key.rk[95]+delta[1][18], 3)
		key.rk[111] = rol(key.rk[103]+delta[2][21], 11)
		key.rk[119] = rol(key.rk[111]+delta[3][24], 17)
		key.rk[127] = rol(key.rk[119]+delta[5][22], 3)
		key.rk[135] = rol(key.rk[127]+delta[6][25], 11)
		key.rk[143] = rol(key.rk[135]+delta[7][28], 17)
		key.rk[151] = rol(key.rk[143]+delta[1][26], 3)
		key.rk[159] = rol(key.rk[151]+delta[2][29], 11)
		key.rk[167] = rol(key.rk[159]+delta[3][0], 17)
		key.rk[175] = rol(key.rk[167]+delta[5][30], 3)
		key.rk[183] = rol(key.rk[175]+delta[6][1], 11)
		key.rk[191] = rol(key.rk[183]+delta[7][4], 17)
		break

	default:
		return
	}

	key.round = (mk_len >> 1) + 16
}

func (key *lea_key) lea_encrypt(ct, pt []byte) {
	X0 := loadU32(pt, 0)
	X1 := loadU32(pt, 1)
	X2 := loadU32(pt, 2)
	X3 := loadU32(pt, 3)

	X3 = ror((X2^key.rk[4])+(X3^key.rk[5]), 3)
	X2 = ror((X1^key.rk[2])+(X2^key.rk[3]), 5)
	X1 = rol((X0^key.rk[0])+(X1^key.rk[1]), 9)
	X0 = ror((X3^key.rk[10])+(X0^key.rk[11]), 3)
	X3 = ror((X2^key.rk[8])+(X3^key.rk[9]), 5)
	X2 = rol((X1^key.rk[6])+(X2^key.rk[7]), 9)
	X1 = ror((X0^key.rk[16])+(X1^key.rk[17]), 3)
	X0 = ror((X3^key.rk[14])+(X0^key.rk[15]), 5)
	X3 = rol((X2^key.rk[12])+(X3^key.rk[13]), 9)
	X2 = ror((X1^key.rk[22])+(X2^key.rk[23]), 3)
	X1 = ror((X0^key.rk[20])+(X1^key.rk[21]), 5)
	X0 = rol((X3^key.rk[18])+(X0^key.rk[19]), 9)

	X3 = ror((X2^key.rk[28])+(X3^key.rk[29]), 3)
	X2 = ror((X1^key.rk[26])+(X2^key.rk[27]), 5)
	X1 = rol((X0^key.rk[24])+(X1^key.rk[25]), 9)
	X0 = ror((X3^key.rk[34])+(X0^key.rk[35]), 3)
	X3 = ror((X2^key.rk[32])+(X3^key.rk[33]), 5)
	X2 = rol((X1^key.rk[30])+(X2^key.rk[31]), 9)
	X1 = ror((X0^key.rk[40])+(X1^key.rk[41]), 3)
	X0 = ror((X3^key.rk[38])+(X0^key.rk[39]), 5)
	X3 = rol((X2^key.rk[36])+(X3^key.rk[37]), 9)
	X2 = ror((X1^key.rk[46])+(X2^key.rk[47]), 3)
	X1 = ror((X0^key.rk[44])+(X1^key.rk[45]), 5)
	X0 = rol((X3^key.rk[42])+(X0^key.rk[43]), 9)

	X3 = ror((X2^key.rk[52])+(X3^key.rk[53]), 3)
	X2 = ror((X1^key.rk[50])+(X2^key.rk[51]), 5)
	X1 = rol((X0^key.rk[48])+(X1^key.rk[49]), 9)
	X0 = ror((X3^key.rk[58])+(X0^key.rk[59]), 3)
	X3 = ror((X2^key.rk[56])+(X3^key.rk[57]), 5)
	X2 = rol((X1^key.rk[54])+(X2^key.rk[55]), 9)
	X1 = ror((X0^key.rk[64])+(X1^key.rk[65]), 3)
	X0 = ror((X3^key.rk[62])+(X0^key.rk[63]), 5)
	X3 = rol((X2^key.rk[60])+(X3^key.rk[61]), 9)
	X2 = ror((X1^key.rk[70])+(X2^key.rk[71]), 3)
	X1 = ror((X0^key.rk[68])+(X1^key.rk[69]), 5)
	X0 = rol((X3^key.rk[66])+(X0^key.rk[67]), 9)

	X3 = ror((X2^key.rk[76])+(X3^key.rk[77]), 3)
	X2 = ror((X1^key.rk[74])+(X2^key.rk[75]), 5)
	X1 = rol((X0^key.rk[72])+(X1^key.rk[73]), 9)
	X0 = ror((X3^key.rk[82])+(X0^key.rk[83]), 3)
	X3 = ror((X2^key.rk[80])+(X3^key.rk[81]), 5)
	X2 = rol((X1^key.rk[78])+(X2^key.rk[79]), 9)
	X1 = ror((X0^key.rk[88])+(X1^key.rk[89]), 3)
	X0 = ror((X3^key.rk[86])+(X0^key.rk[87]), 5)
	X3 = rol((X2^key.rk[84])+(X3^key.rk[85]), 9)
	X2 = ror((X1^key.rk[94])+(X2^key.rk[95]), 3)
	X1 = ror((X0^key.rk[92])+(X1^key.rk[93]), 5)
	X0 = rol((X3^key.rk[90])+(X0^key.rk[91]), 9)

	X3 = ror((X2^key.rk[100])+(X3^key.rk[101]), 3)
	X2 = ror((X1^key.rk[98])+(X2^key.rk[99]), 5)
	X1 = rol((X0^key.rk[96])+(X1^key.rk[97]), 9)
	X0 = ror((X3^key.rk[106])+(X0^key.rk[107]), 3)
	X3 = ror((X2^key.rk[104])+(X3^key.rk[105]), 5)
	X2 = rol((X1^key.rk[102])+(X2^key.rk[103]), 9)
	X1 = ror((X0^key.rk[112])+(X1^key.rk[113]), 3)
	X0 = ror((X3^key.rk[110])+(X0^key.rk[111]), 5)
	X3 = rol((X2^key.rk[108])+(X3^key.rk[109]), 9)
	X2 = ror((X1^key.rk[118])+(X2^key.rk[119]), 3)
	X1 = ror((X0^key.rk[116])+(X1^key.rk[117]), 5)
	X0 = rol((X3^key.rk[114])+(X0^key.rk[115]), 9)

	X3 = ror((X2^key.rk[124])+(X3^key.rk[125]), 3)
	X2 = ror((X1^key.rk[122])+(X2^key.rk[123]), 5)
	X1 = rol((X0^key.rk[120])+(X1^key.rk[121]), 9)
	X0 = ror((X3^key.rk[130])+(X0^key.rk[131]), 3)
	X3 = ror((X2^key.rk[128])+(X3^key.rk[129]), 5)
	X2 = rol((X1^key.rk[126])+(X2^key.rk[127]), 9)
	X1 = ror((X0^key.rk[136])+(X1^key.rk[137]), 3)
	X0 = ror((X3^key.rk[134])+(X0^key.rk[135]), 5)
	X3 = rol((X2^key.rk[132])+(X3^key.rk[133]), 9)
	X2 = ror((X1^key.rk[142])+(X2^key.rk[143]), 3)
	X1 = ror((X0^key.rk[140])+(X1^key.rk[141]), 5)
	X0 = rol((X3^key.rk[138])+(X0^key.rk[139]), 9)

	if key.round > 24 {
		X3 = ror((X2^key.rk[148])+(X3^key.rk[149]), 3)
		X2 = ror((X1^key.rk[146])+(X2^key.rk[147]), 5)
		X1 = rol((X0^key.rk[144])+(X1^key.rk[145]), 9)
		X0 = ror((X3^key.rk[154])+(X0^key.rk[155]), 3)
		X3 = ror((X2^key.rk[152])+(X3^key.rk[153]), 5)
		X2 = rol((X1^key.rk[150])+(X2^key.rk[151]), 9)
		X1 = ror((X0^key.rk[160])+(X1^key.rk[161]), 3)
		X0 = ror((X3^key.rk[158])+(X0^key.rk[159]), 5)
		X3 = rol((X2^key.rk[156])+(X3^key.rk[157]), 9)
		X2 = ror((X1^key.rk[166])+(X2^key.rk[167]), 3)
		X1 = ror((X0^key.rk[164])+(X1^key.rk[165]), 5)
		X0 = rol((X3^key.rk[162])+(X0^key.rk[163]), 9)
	}

	if key.round > 28 {
		X3 = ror((X2^key.rk[172])+(X3^key.rk[173]), 3)
		X2 = ror((X1^key.rk[170])+(X2^key.rk[171]), 5)
		X1 = rol((X0^key.rk[168])+(X1^key.rk[169]), 9)
		X0 = ror((X3^key.rk[178])+(X0^key.rk[179]), 3)
		X3 = ror((X2^key.rk[176])+(X3^key.rk[177]), 5)
		X2 = rol((X1^key.rk[174])+(X2^key.rk[175]), 9)
		X1 = ror((X0^key.rk[184])+(X1^key.rk[185]), 3)
		X0 = ror((X3^key.rk[182])+(X0^key.rk[183]), 5)
		X3 = rol((X2^key.rk[180])+(X3^key.rk[181]), 9)
		X2 = ror((X1^key.rk[190])+(X2^key.rk[191]), 3)
		X1 = ror((X0^key.rk[188])+(X1^key.rk[189]), 5)
		X0 = rol((X3^key.rk[186])+(X0^key.rk[187]), 9)
	}

	/**
	_ct[0] = loadU32(X0);
	_ct[1] = loadU32(X1);
	_ct[2] = loadU32(X2);
	_ct[3] = loadU32(X3);
	*/

	ct[0] = byte(X0)
	ct[1] = byte(X0 >> 8)
	ct[2] = byte(X0 >> 16)
	ct[3] = byte(X0 >> 24)

	ct[4] = byte(X1)
	ct[5] = byte(X1 >> 8)
	ct[6] = byte(X1 >> 16)
	ct[7] = byte(X1 >> 24)

	ct[8] = byte(X2)
	ct[9] = byte(X2 >> 8)
	ct[10] = byte(X2 >> 16)
	ct[11] = byte(X2 >> 24)

	ct[12] = byte(X3)
	ct[13] = byte(X3 >> 8)
	ct[14] = byte(X3 >> 16)
	ct[15] = byte(X3 >> 24)
}

func (key *lea_key) lea_decrypt(pt, ct []byte) {
	X0 := loadU32(ct, 0)
	X1 := loadU32(ct, 1)
	X2 := loadU32(ct, 2)
	X3 := loadU32(ct, 3)

	if key.round > 28 {
		X0 = (ror(X0, 9) - (X3 ^ key.rk[186])) ^ key.rk[187]
		X1 = (rol(X1, 5) - (X0 ^ key.rk[188])) ^ key.rk[189]
		X2 = (rol(X2, 3) - (X1 ^ key.rk[190])) ^ key.rk[191]
		X3 = (ror(X3, 9) - (X2 ^ key.rk[180])) ^ key.rk[181]
		X0 = (rol(X0, 5) - (X3 ^ key.rk[182])) ^ key.rk[183]
		X1 = (rol(X1, 3) - (X0 ^ key.rk[184])) ^ key.rk[185]
		X2 = (ror(X2, 9) - (X1 ^ key.rk[174])) ^ key.rk[175]
		X3 = (rol(X3, 5) - (X2 ^ key.rk[176])) ^ key.rk[177]
		X0 = (rol(X0, 3) - (X3 ^ key.rk[178])) ^ key.rk[179]
		X1 = (ror(X1, 9) - (X0 ^ key.rk[168])) ^ key.rk[169]
		X2 = (rol(X2, 5) - (X1 ^ key.rk[170])) ^ key.rk[171]
		X3 = (rol(X3, 3) - (X2 ^ key.rk[172])) ^ key.rk[173]
	}

	if key.round > 24 {
		X0 = (ror(X0, 9) - (X3 ^ key.rk[162])) ^ key.rk[163]
		X1 = (rol(X1, 5) - (X0 ^ key.rk[164])) ^ key.rk[165]
		X2 = (rol(X2, 3) - (X1 ^ key.rk[166])) ^ key.rk[167]
		X3 = (ror(X3, 9) - (X2 ^ key.rk[156])) ^ key.rk[157]
		X0 = (rol(X0, 5) - (X3 ^ key.rk[158])) ^ key.rk[159]
		X1 = (rol(X1, 3) - (X0 ^ key.rk[160])) ^ key.rk[161]
		X2 = (ror(X2, 9) - (X1 ^ key.rk[150])) ^ key.rk[151]
		X3 = (rol(X3, 5) - (X2 ^ key.rk[152])) ^ key.rk[153]
		X0 = (rol(X0, 3) - (X3 ^ key.rk[154])) ^ key.rk[155]
		X1 = (ror(X1, 9) - (X0 ^ key.rk[144])) ^ key.rk[145]
		X2 = (rol(X2, 5) - (X1 ^ key.rk[146])) ^ key.rk[147]
		X3 = (rol(X3, 3) - (X2 ^ key.rk[148])) ^ key.rk[149]
	}

	X0 = (ror(X0, 9) - (X3 ^ key.rk[138])) ^ key.rk[139]
	X1 = (rol(X1, 5) - (X0 ^ key.rk[140])) ^ key.rk[141]
	X2 = (rol(X2, 3) - (X1 ^ key.rk[142])) ^ key.rk[143]
	X3 = (ror(X3, 9) - (X2 ^ key.rk[132])) ^ key.rk[133]
	X0 = (rol(X0, 5) - (X3 ^ key.rk[134])) ^ key.rk[135]
	X1 = (rol(X1, 3) - (X0 ^ key.rk[136])) ^ key.rk[137]
	X2 = (ror(X2, 9) - (X1 ^ key.rk[126])) ^ key.rk[127]
	X3 = (rol(X3, 5) - (X2 ^ key.rk[128])) ^ key.rk[129]
	X0 = (rol(X0, 3) - (X3 ^ key.rk[130])) ^ key.rk[131]
	X1 = (ror(X1, 9) - (X0 ^ key.rk[120])) ^ key.rk[121]
	X2 = (rol(X2, 5) - (X1 ^ key.rk[122])) ^ key.rk[123]
	X3 = (rol(X3, 3) - (X2 ^ key.rk[124])) ^ key.rk[125]

	X0 = (ror(X0, 9) - (X3 ^ key.rk[114])) ^ key.rk[115]
	X1 = (rol(X1, 5) - (X0 ^ key.rk[116])) ^ key.rk[117]
	X2 = (rol(X2, 3) - (X1 ^ key.rk[118])) ^ key.rk[119]
	X3 = (ror(X3, 9) - (X2 ^ key.rk[108])) ^ key.rk[109]
	X0 = (rol(X0, 5) - (X3 ^ key.rk[110])) ^ key.rk[111]
	X1 = (rol(X1, 3) - (X0 ^ key.rk[112])) ^ key.rk[113]
	X2 = (ror(X2, 9) - (X1 ^ key.rk[102])) ^ key.rk[103]
	X3 = (rol(X3, 5) - (X2 ^ key.rk[104])) ^ key.rk[105]
	X0 = (rol(X0, 3) - (X3 ^ key.rk[106])) ^ key.rk[107]
	X1 = (ror(X1, 9) - (X0 ^ key.rk[96])) ^ key.rk[97]
	X2 = (rol(X2, 5) - (X1 ^ key.rk[98])) ^ key.rk[99]
	X3 = (rol(X3, 3) - (X2 ^ key.rk[100])) ^ key.rk[101]

	X0 = (ror(X0, 9) - (X3 ^ key.rk[90])) ^ key.rk[91]
	X1 = (rol(X1, 5) - (X0 ^ key.rk[92])) ^ key.rk[93]
	X2 = (rol(X2, 3) - (X1 ^ key.rk[94])) ^ key.rk[95]
	X3 = (ror(X3, 9) - (X2 ^ key.rk[84])) ^ key.rk[85]
	X0 = (rol(X0, 5) - (X3 ^ key.rk[86])) ^ key.rk[87]
	X1 = (rol(X1, 3) - (X0 ^ key.rk[88])) ^ key.rk[89]
	X2 = (ror(X2, 9) - (X1 ^ key.rk[78])) ^ key.rk[79]
	X3 = (rol(X3, 5) - (X2 ^ key.rk[80])) ^ key.rk[81]
	X0 = (rol(X0, 3) - (X3 ^ key.rk[82])) ^ key.rk[83]
	X1 = (ror(X1, 9) - (X0 ^ key.rk[72])) ^ key.rk[73]
	X2 = (rol(X2, 5) - (X1 ^ key.rk[74])) ^ key.rk[75]
	X3 = (rol(X3, 3) - (X2 ^ key.rk[76])) ^ key.rk[77]

	X0 = (ror(X0, 9) - (X3 ^ key.rk[66])) ^ key.rk[67]
	X1 = (rol(X1, 5) - (X0 ^ key.rk[68])) ^ key.rk[69]
	X2 = (rol(X2, 3) - (X1 ^ key.rk[70])) ^ key.rk[71]
	X3 = (ror(X3, 9) - (X2 ^ key.rk[60])) ^ key.rk[61]
	X0 = (rol(X0, 5) - (X3 ^ key.rk[62])) ^ key.rk[63]
	X1 = (rol(X1, 3) - (X0 ^ key.rk[64])) ^ key.rk[65]
	X2 = (ror(X2, 9) - (X1 ^ key.rk[54])) ^ key.rk[55]
	X3 = (rol(X3, 5) - (X2 ^ key.rk[56])) ^ key.rk[57]
	X0 = (rol(X0, 3) - (X3 ^ key.rk[58])) ^ key.rk[59]
	X1 = (ror(X1, 9) - (X0 ^ key.rk[48])) ^ key.rk[49]
	X2 = (rol(X2, 5) - (X1 ^ key.rk[50])) ^ key.rk[51]
	X3 = (rol(X3, 3) - (X2 ^ key.rk[52])) ^ key.rk[53]

	X0 = (ror(X0, 9) - (X3 ^ key.rk[42])) ^ key.rk[43]
	X1 = (rol(X1, 5) - (X0 ^ key.rk[44])) ^ key.rk[45]
	X2 = (rol(X2, 3) - (X1 ^ key.rk[46])) ^ key.rk[47]
	X3 = (ror(X3, 9) - (X2 ^ key.rk[36])) ^ key.rk[37]
	X0 = (rol(X0, 5) - (X3 ^ key.rk[38])) ^ key.rk[39]
	X1 = (rol(X1, 3) - (X0 ^ key.rk[40])) ^ key.rk[41]
	X2 = (ror(X2, 9) - (X1 ^ key.rk[30])) ^ key.rk[31]
	X3 = (rol(X3, 5) - (X2 ^ key.rk[32])) ^ key.rk[33]
	X0 = (rol(X0, 3) - (X3 ^ key.rk[34])) ^ key.rk[35]
	X1 = (ror(X1, 9) - (X0 ^ key.rk[24])) ^ key.rk[25]
	X2 = (rol(X2, 5) - (X1 ^ key.rk[26])) ^ key.rk[27]
	X3 = (rol(X3, 3) - (X2 ^ key.rk[28])) ^ key.rk[29]

	X0 = (ror(X0, 9) - (X3 ^ key.rk[18])) ^ key.rk[19]
	X1 = (rol(X1, 5) - (X0 ^ key.rk[20])) ^ key.rk[21]
	X2 = (rol(X2, 3) - (X1 ^ key.rk[22])) ^ key.rk[23]
	X3 = (ror(X3, 9) - (X2 ^ key.rk[12])) ^ key.rk[13]
	X0 = (rol(X0, 5) - (X3 ^ key.rk[14])) ^ key.rk[15]
	X1 = (rol(X1, 3) - (X0 ^ key.rk[16])) ^ key.rk[17]
	X2 = (ror(X2, 9) - (X1 ^ key.rk[6])) ^ key.rk[7]
	X3 = (rol(X3, 5) - (X2 ^ key.rk[8])) ^ key.rk[9]
	X0 = (rol(X0, 3) - (X3 ^ key.rk[10])) ^ key.rk[11]
	X1 = (ror(X1, 9) - (X0 ^ key.rk[0])) ^ key.rk[1]
	X2 = (rol(X2, 5) - (X1 ^ key.rk[2])) ^ key.rk[3]
	X3 = (rol(X3, 3) - (X2 ^ key.rk[4])) ^ key.rk[5]

	/**
	_pt[0] = loadU32(X0)
	_pt[1] = loadU32(X1)
	_pt[2] = loadU32(X2)
	_pt[3] = loadU32(X3)
	*/

	pt[0] = byte(X0)
	pt[1] = byte(X0 >> 8)
	pt[2] = byte(X0 >> 16)
	pt[3] = byte(X0 >> 24)

	pt[4] = byte(X1)
	pt[5] = byte(X1 >> 8)
	pt[6] = byte(X1 >> 16)
	pt[7] = byte(X1 >> 24)

	pt[8] = byte(X2)
	pt[9] = byte(X2 >> 8)
	pt[10] = byte(X2 >> 16)
	pt[11] = byte(X2 >> 24)

	pt[12] = byte(X3)
	pt[13] = byte(X3 >> 8)
	pt[14] = byte(X3 >> 16)
	pt[15] = byte(X3 >> 24)
}

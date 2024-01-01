package lsh256

import (
	. "github.com/RyuaNerin/testingutil"
)

var testCases224 = []HashTestCase{
	// 암호알고리즘 검증기준 V3.0
	// 테스트 벡터
	// LSH(256-224)ShortMsg.txt
	{
		Msg: ``,
		MD:  `48A0D55B2B3D91F26E06F7110FE9CE8EA0E2656BBE344CB1C5930653`,
	},
	{
		Msg: `56`,
		MD:  `6B3F5B095692AF487964F51B3248FE88132AD3F428969102AB3C2CF3`,
	},
	{
		Msg: `6C74`,
		MD:  `BA98E22D28A98DEF2A1D795C39D50758956433791DAE8FBDCCEF4312`,
	},
	{
		Msg: `B0756E`,
		MD:  `D9ADFC16779DF16D3099012D568FABA4637A63414F2A970323B9154B`,
	},
	{
		Msg: `9CF4B3D5`,
		MD:  `AB33EA8BBE92193C7B2B520FBEB9D75C0F87D4FF9DCF082BBD18EA30`,
	},
	{
		Msg: `43A8DAD425`,
		MD:  `51E392DBF97F8909936BCEB4B4087B9F2EEA968F449BBA1A7FE92BFE`,
	},
	{
		Msg: `AA215AC188AC`,
		MD:  `E0FE16C6F094BF3515EF2D4F1045F506073185412EB41157525E3B28`,
	},
	{
		Msg: `890395536B1F3C`,
		MD:  `C7CFEDA470588FA2C0FD19ADAC3155277E17BF783D400171F8E9479B`,
	},
	{
		Msg: `A4E0C0952BDB6236`,
		MD:  `DEA30B8F052318E676EB52BFDE7E2F9B02D132E49F1F27ABFC8ADB29`,
	},
	{
		Msg: `F1CEFC4C4263E3EAB3`,
		MD:  `A2CFAB024343896AA4144D789A93089003A1CFE3A44F2A752BA809D3`,
	},
	{
		Msg: `6DEBA0641ECF50837AD3`,
		MD:  `DF63AB70AACBD6EC123DCCDD4013B1FC783096E135CD37EF25E304B0`,
	},
	{
		Msg: `3D110BA765E9D78F11FE8F`,
		MD:  `D9EE59A60D857A95D38AB5E99DB47FB09ED64CDF4C897CBBF6278011`,
	},
	{
		Msg: `C1AAE4090ECE92F0C5A366BB`,
		MD:  `3CB0B91BBE2E99F89BF30BBF50247B3CCE68B17DFBF05B36C0C390FE`,
	},
	{
		Msg: `D56B69DB3159D53979A3D4DFD6`,
		MD:  `5D4705338929C23F2D1766B7326C48A8AA0553A423F0EF04BFD8C8C6`,
	},
	{
		Msg: `D98BF5267557F059B356061AA32A`,
		MD:  `219B7C07E181674392DB2C131D183C12BFC4B91B7BB4642F63B481DC`,
	},
	{
		Msg: `E5B27CF9A42A33B8FAB04CAAAE3DF1`,
		MD:  `F002DA8624F939E9F66DFA2B160DD370899D7991B7B270FF6CFA4C21`,
	},
	{
		Msg: `234E74419224CCEABE66F1FA50AD68F2`,
		MD:  `1DAB9A667311CC2D176D17271905B2F722E36D42A808BC8265896CCB`,
	},
	{
		Msg: `92A8721C7612341FDAE114D10AC1664C87`,
		MD:  `6FCDBC9EB4D342548BAE75992A3EFCCA819E39A85AA43719420002FE`,
	},
	{
		Msg: `C04844074BCA98BC595326D9E640E453A3CD`,
		MD:  `61742E9F2CB6CD25D15BD59D6990F4F320EA6488832E3D77A6D6FD3E`,
	},
	{
		Msg: `2DE4EC395856FA258FFE5A3EDE879D1B8BF152`,
		MD:  `C3B4096DF31C6BFB3A9BD565AFEAF32211FE4BFED342724700D169CA`,
	},
	{
		Msg: `68A6D83145E0FA99E9A1F43D8B544C1A8206B4AC`,
		MD:  `7A98171FFF3114C438F7A113393433DAB1BF59F7B139ACE474B9F26F`,
	},
	{
		Msg: `AB2F8D610BCFAB4EB5BD5820EE3C9B5690D0BAAEAC`,
		MD:  `04FAE07987F9DAE38E00E2AE9D02CB101E152BB74FA8C4FF2A5F52BF`,
	},
	{
		Msg: `50FDDD733BB583574F7D5E6D58C0BE16477B39DE8344`,
		MD:  `CB41D048C243B5D79908357748DA26C0F7E01C14AAE0361F6985063C`,
	},
	{
		Msg: `9FA4848C31FC9D210DD119F6ED4F8A2711FF5E387A34B3`,
		MD:  `789BF6B08427F103778F06A1BB07DD2491B0E0D831ABA23C59A1575A`,
	},
	{
		Msg: `9E59A76D80C75E85BEB6D2A3332850F0379FA02568D4E5E9`,
		MD:  `B50694EE0C24D3E7C0C9046B57F32D017B5C51012BC9864204A71593`,
	},
	{
		Msg: `ABABC4596D67C8E82F1A05694160FC6ACD4DBFC418A4567005`,
		MD:  `4CB5ECC8C7EFA3E011526C0CE513CAC6725BEAA720F38C5D5E914926`,
	},
	{
		Msg: `4A609CDBFAFDFE4739EC69A58DF19E7BEA7A64466D81DB7A4DB4`,
		MD:  `1FFE7816ABAB4B1E3E89C0CFE9F2CA917ABF4F813BC8522721865040`,
	},
	{
		Msg: `0DC40ED4C692F75BA6D5C7617405A4C0B56654777D8BAF4C861854`,
		MD:  `41D895E3933C62B37B24793C29CF0A3AB19559787F71A5F031603B9F`,
	},
	{
		Msg: `B67FCC9CF9A6BECC9399E68FD769E0A3D8F1EDD288663B24212F6653`,
		MD:  `F70CD9B6D2A53C6173AADED5F8F2CA550C7B14D2491991A82C556BD5`,
	},
	{
		Msg: `C4F63E7F60FFCB23D6B88E485B920AF81D1083F6291D06AC8CA3A965B8`,
		MD:  `61141F7334F09AFCEA38CA6270E9B6C48DAA702AB7C5F273005613FA`,
	},
	{
		Msg: `5914BC2ADD40544A027FCA936BBDE8F359051C75C147281DA2581A13953A`,
		MD:  `7DB20BF0FBFF14F5A9E6111D51FF15F5C6B77689308B6508F9B3F1F5`,
	},
	{
		Msg: `7F05A69ED4C662A73D3B31F17A77E6D2162595DFE126A4B2A17C25EB056763`,
		MD:  `9EFF1FC6015D8D13C3B249E0C990F99A7CD800EE570392D1D3412F1E`,
	},
	{
		Msg: `75DACDBB0B81DFAEFAA1FA342819CB525539B7EEA1317AE34EC90E51C853CE36`,
		MD:  `1CF19B08016C1B26615D24270823234148BFF61B527BC5AE9464AA4C`,
	},
	{
		Msg: `E2DEDB4363ADBFE8C1E850B22C90D7D6CEB8720E52F3B7A3A3099ACBF6EE91360A`,
		MD:  `1AEA263027F3CDB9E24AE2208BCBF6BBDC6E33BF53DEA9E5871416C6`,
	},
	{
		Msg: `F4F6E4AE3FAC398579055306C8E39BC6FD1E9AFA0F260F37C7EA3E7DBCF88A7482E5`,
		MD:  `D0E0CB9903D8BB3B67C360126A1B6EF1B84D03C68217F0C4E7756D23`,
	},
	{
		Msg: `46F51FEFC9CCBEE1C0BDA837C9C56F833DA16C0F4C956CC45066A74A5BA8F271124035`,
		MD:  `4DAC71C7526C07032745E16506997667F4FC35400B488BC1B4492DBF`,
	},
	{
		Msg: `D0BF211EACA159171B875176689D50F4BA854261F1A9B896D1D02FAA62F0956996430901`,
		MD:  `41094A3D97E3426F4B45F5B998D7F566168EE7F48F8077481B49E949`,
	},
	{
		Msg: `2D9DE3FFD6FDF47A97620F8E6EE0C70D39F27B70D31290148087FEA38B7F0ECF2D5D6247B1`,
		MD:  `6AA9AED39F45241DF283E0478F66D8BAC4A29CD8B2803EAABDEF90AE`,
	},
	{
		Msg: `5435E3D635E4227585ADD763D3B71F99485396928082CBF598B7F17489DED7F4B49C1AEEBED9`,
		MD:  `BA66EF69F32F8F6929572F800FC127BF17F0914A271B9DBAAA4BB62F`,
	},
	{
		Msg: `99159D7C2939F3758A6B9B019A628229C91DB98CCDCAC85F97E744B285A7A4F5D99750BA6A81FB`,
		MD:  `721DFF4416AD05DAD4142B35A4A5E30E4C05164EA504D2F5499CAC3D`,
	},
	{
		Msg: `BD9990332187B636FC854B0816CCF7C2EE6FD12DFE8DC88F40CB43C5F1010492953452373BC55744`,
		MD:  `548AF5945579C8BABD8C085D955691BAC94F426368E1561568BDA3F7`,
	},
	{
		Msg: `02E7C7ADDC56C2209247FC96099A1B29C3576025A8588A610743F7D0D8BC0D5344D8D8D167D4F48A54`,
		MD:  `B43E1B02E72206C893537CB816CEB75FA216B960EC2CF2E5FCA5CD6A`,
	},
	{
		Msg: `07AEA8591566178566389768D6DB101AAEC8257FD49F0BC83EC79C3DAC55A73A97308B786472048EDE11`,
		MD:  `ADC93CBBC5CC7B0517D3734B17D3783698ABFB434628BDE584FEEACD`,
	},
	{
		Msg: `44B8E3BE46D65D0D5FB875D7169803DD35E848AD204FDFCA66EEF10796E7A9345420E7FCAB44390EA50B07`,
		MD:  `81C743B8E9CD8FA82C870807148D930E352F785DA98D23F2B6F28523`,
	},
	{
		Msg: `7B19BC0D5FC5F00EC49C98C7C3614123EC7597ADE9EE284A5DE86B28122001F254BF3A469873DB6862C41F38`,
		MD:  `3DC35321AB72E231D6738D532ACBA6D4C880079B5F9DAA35CAE9B83E`,
	},
	{
		Msg: `65A6FB61EB4AA6D4773A5975A9A5541363F7219E2AAD1C94CF6833BD4E40EB18F65AA315CABD1D0428613ACEA8`,
		MD:  `2CB370CB5FA32E60D1466D587836F9AF8DC96887991BAE0C536ABD43`,
	},
	{
		Msg: `8905E34753F373C58E6AE14BCB7F7551B70451BFA8C2A652750327428B54BE223D887E075C00DAEC082A1455C548`,
		MD:  `E955C4AC5D0FE789EC1EE7AED45FD62278341BCE3B52E53DDAD82B92`,
	},
	{
		Msg: `F833D273C908B6DF7C4B1B60240EDF2DAEDE69E5CEEB39CF3969C36615357E15ADBC383F6A19AC82FF7D563DE3D486`,
		MD:  `5785087F1D619CFA7B008BE7F1B8760BF627FD9457278E9B1A4E2026`,
	},
	{
		Msg: `4F04A5B422BA3A13E9F9026825CC331D6EFEB391701A39B33497B0CFA14CDCA8422A3FC649293E11D483809F702197CB`,
		MD:  `B889E1965EF80A3A659915452DBD8C6CA4CFE7B136E749E39F6A748C`,
	},
	{
		Msg: `DB52E2870FC3B1F37BACAB2E94476D6B8C86744861053D5C1A3EAA0F4EFCA8AFE5120FC0224CCA93E2A6371A1EBCD0FAB7`,
		MD:  `DFA2CC9CF02DA1705593F52957CA4202A0CEC33D94867942666F123A`,
	},
	{
		Msg: `820CE27335801B42914A61C745AFF213A403E4B49DC9BBEEED1C3B48934E7B4B84CE3B86A14F77D5E6263DE3E4F4BBEE7186`,
		MD:  `0F91307A6B383A55FA97B1B27A6ACDFC016C380D19FAE98699F1C199`,
	},
	{
		Msg: `9CCA9EF99529E9223A19B42D17CEF683FCA231F2985AA40849ABC1BF33FDE160906E37A1461A0210BF9A07B5F5CA3FF3CCCFCD`,
		MD:  `1A90ACBABA9E26562BE66BFDA71F0F2EE973DC7FA5F4D51E0BD88F19`,
	},
	{
		Msg: `B0A55932B63026C45881C85D5867A75008F949A1B071D175303EC6C7E654165DD2C75BFAF8985CD606CFD74B37C97ED177D43F2C`,
		MD:  `1D148C51E13D81C6758C36BD7B3ED272AFD7A8C691E445825AA15AAF`,
	},
	{
		Msg: `C8F9C8E4E4194BF6243C7EC9A025E763982A56DD1342B314C560F9A236FF830D21FE873D1AAFE034742B10F962E885D55C8D1BC2C5`,
		MD:  `CF1524D7EED9BD4014C32EC96C21B2F659640B199A72E68701CB941B`,
	},
	{
		Msg: `04E321CAD123AB8DA39DAA832914691A270DE9A0C9D3D56037BDB80E69B36E653DEFD90EC701D5C4E8166DB622794F2CAA5A7D6EBFCF`,
		MD:  `87D69E9BE5A3B3D06EFE9817AAB6D29DB4F9D41A4A27938B6F6820A2`,
	},
	{
		Msg: `46E23B1875746827C83A1460A358E89D0622F575BD8190BED698605058A7B50C322F1B1D7257307EB78EB2BCE972DCFE27F0F34201AD1C`,
		MD:  `5E0D8C26DF888801CDF74FAEFD516D3EAD327F5FFACE95FA9ACAD057`,
	},
	{
		Msg: `E9E8A57076E213529383900055187E69C025A73F8E12224C2C21770906A3026D046C9D035484793F6B1F4D1246CFB3BCF6F9CEE25B8631E0`,
		MD:  `0E897FA1F9129E85DCDEB7823BFD2B2D1B1FB1008364B85DD158AAD6`,
	},
	{
		Msg: `30DAFC60470A482C1649C64073A4A2249CF0AE787A00BCD0271B4BC3ADAC673F113CCA1F909E4700286969505AAC57A72FDCFA5B7911624F5B`,
		MD:  `3BF05570599676E4107C557966C92ECA8AA8558B7377E8359FB8C024`,
	},
	{
		Msg: `72660B0BEE8B12D72E286F66C7C79241B3CD5A91025E732BCD7BE2C96C6DDE9AE867E2FEB367008A975D7853ED8F89690F3C87A1107F2E98AA77`,
		MD:  `8A4B2BC476ECE69E862E151508F2176C2F499D2DF1A89BC436430494`,
	},
	{
		Msg: `36F477A527ED64956F0D64C1B23361B261DE78688EA865FCFF113C84817E5B377E829CD2D25BCF3ADBC06762CFDA736F5390D01A49079D56E969F0`,
		MD:  `C453E8E5E51CE06B895600DAD258A6A0A62133DEB6D3727F9BD68B0D`,
	},
	{
		Msg: `3313E6C703E3F942BB87ED0F9C4D9F25120085B5DC75EF5D6D618DA0926D3293568DD7D8238DE3D029EDAA0D03EE330698ED3E0F40F08588297F4354`,
		MD:  `F159A9EC6919560E2FB0D491E51250251305571E93AC95ACEFB7A1D5`,
	},
	{
		Msg: `F8B5F38DB3965A75DEA3C7B0E4AC9AB1B9E60373BE3D0C69E5EBDAB7154A2FA2D0E6344513D776867D4E2E49B46DC5E4055C97D24A99F8AB118F7453FA`,
		MD:  `8E2843213DCD9CF7A6021C81156DDDE0613F171F0C8EF568F21E94BA`,
	},
	{
		Msg: `6B5B72612B7E253DF733F448B59F0361D8348C2E78788C05A808BD0D04DDC5A089CFADF4AA5794FB9F99A58701A246B294EE927D642D0AB713468721719D`,
		MD:  `3318B375C603E1825C4F92C17A5BC99D4747E57E50AFF890BBA7E166`,
	},
	{
		Msg: `74CF4D9239D18B4AFA1377B1828259BF7A704AFBDF3A453EB32881BAF294FE3AC74F5EAB0040F7DCC04E26827381E67F0763A1645AAD4634EE1F9B6BE4B0D5`,
		MD:  `3BAB8D35F8A527B5C3DBDF053974050E8B0A8F53D1E5597CA63E1B49`,
	},
	{
		Msg: `3E28E76625E98916974C74551461CBCF6956534230EA1D8C962566096146BB1B31FE610A05DBF30F02074DAB7032A93D06CD08D46779D442ECBFC372A401FEFB`,
		MD:  `F21649E2E5C93B954AA8D6305F4ED391529A481B452E326D94C058AA`,
	},
	{
		Msg: `99B083F455C0C70BA559A0D19DD06358DCD5490012EE188EBF24FF7E3DFBC19B00A816CC1D7F4E3109AB376EA321414775D7E06E3730329895BBE2D6331E4C23ED`,
		MD:  `0CC1B8B7EBCF33146B6DB460AB7599B1213957DDBE32955C613BDE19`,
	},
	{
		Msg: `4DE1396300CDAAB563D82B8B0F8A33593993C5E1296A89F9EE33228E51E6BCE907AEE4312B8F9D31EB4CB05C8188430E652BAD15C007897024BC09567BD88E7B3E43`,
		MD:  `3E0F6B91A8FFDA8EDF913EE0C07EB6CB6F1C83B10CBB1402CE6F3C9F`,
	},
	{
		Msg: `758EE8DB3306A95BA41E5E81A01EC170D0DDDD53BF84C6433494866BC02B5A6993831FF893BB07CCB0DA90BE7111B729D340689454BC3A914B5059E47A81C3E8E616B8`,
		MD:  `C106181883AB7E23B31186BEF4D3F3682A9A81224A0704919E51636E`,
	},
	{
		Msg: `2FE7D86B697B97FE87DF6DAD2061474A1370A0E96EF590DC8EE606E28817519BC854E63FF8AE5D902412BC0D1B20F8691839B93BA1DA5EF7B029411D11CE3CB00B046CEA`,
		MD:  `8C225C3C399C73C8FFDF8AF142B84D28A6787B84CE21A02EE960FD13`,
	},
	{
		Msg: `42F0342EDC31853EE3B85DAEEADA8E99F4A3C192916EBADD4AB2DEE03F1FF2D40FEA25C0F86551A840D27E01D379A8692DF8C2957ADBFA8BFC6C3F11AECC00652743372EC0`,
		MD:  `43B4D198684126F94488C9B3C4F23066338B6F5B3D66EF4F717CABE0`,
	},
	{
		Msg: `D5A7B3F916614B193F8D14158067E8B398888DEFA4672C9E341F0551A473607CAA0FC7D82B10720553C70374E4BA7484CEEC4DEB015F1A587FD61C57B5DF3F740B114CA3BA93`,
		MD:  `108C505C242009182C9AE48283F6E0BD0D48D0DE58D45E33D601A141`,
	},
	{
		Msg: `60DFDDDE9D12CCC31F35A5AF31B82B2EC124D4AE12FF1D86B1F02136AA236AF0FBEF075849B6DE89D20D4458A31DCBDC0B495B2E4F12856D9F3713054920E4312D7CF7BF712AC5`,
		MD:  `9BE757CA13B9120D5689444D7E16516ECA2D4AB409A48349FA66A157`,
	},
	{
		Msg: `147655D46F0AE48972DF113EDE3B57949EF2920A079262CF5EDDD0D59A7C34276CF1567C7DCED30CF724C2A2463F98F32090492EC2EBB0D47DF331CA72D29DAD9A2D55650956B3FD`,
		MD:  `194E63F7E0A5951BC4ED766C9335210D9957D1960F2D2E47CD795459`,
	},
	{
		Msg: `73C4A5B1EFC19B49FB63A4A6C0BDD39B9856E8FECEDB2CFEA5109DB69B0C4B03EEA60293602B293F3752EA9B897194CBB6F03D9836FEABC395847598BB5765C771B2217021E30A7DDD`,
		MD:  `08ADF7B832F37EE3A6EFFAF1B8AA90072A9D0F9923B257FD9A649271`,
	},
	{
		Msg: `446534A36DB90CE0A877C032A6E0A372702433F29FC7F003D0496FBB5306AAF5235B9A981B0A081D66E0A04B38954745BBF3AEC74D3A27ADE8C4B8B5E1E56E85D615F4F0E251647512BE`,
		MD:  `C4F129B645EED13C1881119D39C8260776835223D552D51ED479B0B8`,
	},
	{
		Msg: `5D081D9056164F2A32C3C7861476E2D87497F76E739745DB6B312954CBF0955A22F2BD9A79EB1311E66C76735E66A6CECBEAFC6A5E3258F527E55ADCE397D5FC158F807A58CEEFEFF34DC2`,
		MD:  `E0BB693A3B8168A0C0F2640CF36351095BA496043DA8F2AB4CC02E20`,
	},
	{
		Msg: `C9AE851DEEE43458A2DA5FF7AA4F93DCB15471BBEC5E6ACFFA9AC6CB0E80EB144FABDC2D7FE80072DB8D250A430CC4179BC7862F71260810E0E7DB8FB4B6D0D2CEFE72BFCBFED3B4E01F8C79`,
		MD:  `E99EEFE399E18C5F7FD8A433E74A2DD2BACA9A5051D197718C9D0223`,
	},
	{
		Msg: `0885FF7362E4DF1AD495F3ADF5E916242ABD83660F4E4FBC184F378578399C2C41E3CC58CD5283625ED2B20C01457E70B8BC8C36D7447C0AAF03B34D6A835E0B364DB32EA6BA7AF15406CD97CC`,
		MD:  `649B75062AD734A6D4696AFA031BB47084A980745174161F3309FF2B`,
	},
	{
		Msg: `A61C5CA5DB74ACF1880093D841D0757FA4D9FF2FEC3D5DB3C3CDAD0158CCB41DB2F6D03B44C12E24C1E33DDA3114637F4515EE73E7F14DBD168380ACF138CDBD1117B5C01CFE2637972EF2FFFC09`,
		MD:  `03511B321CF388E764508242BB4ACC17B625DE939D79D7D9FBC3CFE7`,
	},
	{
		Msg: `80320D1BC24B6758668B5AA4519CF89949756C84F17D8EB0FBF3B4E8998DCA37345FA5F261FF94F811E73839071A9675900B6B466EBABF611520BD191ACDC225E8D0B4A2B5D7955CB49758498C6EF3`,
		MD:  `CA4413AE30A160F537D18344883D919648687013C8FB77F1DEE4930A`,
	},
	{
		Msg: `560A47B56612BFDA98DC865DC95D070D5B024BE334A41A9309B96461DCC7E036C31FF3F3ECC7D305DEB8EC3A014983D261EC99DAAF1812B5898616FCFAA0B8884537D5CCDE743F68BA111BFCC8C8883A`,
		MD:  `D5A7133B59F9A115BA7FA816B6A818FADCCE540A7381E03EBA27C10A`,
	},
	{
		Msg: `C00C82401BF80524628575D0ACFF3D1063BAB7E946014C2DB44FDC4C27C25F6DFFCCCDBD812C3C071F1C4886F9CCB8F7A5105FF12B553F7954DAF5D2251D5192A69344FBE89BFF8745B559C5EAEF716493`,
		MD:  `A226C91EF1EDA5F5A6C0F23521F054914383FBB17F243DF39D24B6D8`,
	},
	{
		Msg: `A1958D965A6377231A2E35CAC2A4D23DC8AD9F582F96CC5EB76E368651B080B6D523658EF7FF4CA878511D1E99A0544DD5D154D9CF48FDF28A5080D5D8BD729515D0201B145A236B5A24342A8A8F415D2178`,
		MD:  `C1D0F6627E4E201CD8809A033DE9B60A541CFC3F918F7498C5ABA031`,
	},
	{
		Msg: `C074635E438796C74CDF6DDF63B8C76367459ABD22F529711FF57EB7C0F614698B30D420C5A7429B42ADF9607BA855A5E1156F7071CD4EFE89875B62FA9F166E5E1608D6ACD9F89309466286E1EC255E37927C`,
		MD:  `55E8EAEEEA5013D0DB75E518B2AE0AC777E44E42790105A1C735C886`,
	},
	{
		Msg: `BB397B2F9318E90BD616E08F4DE12CE8C81FA620784AB2929ABE0CA15510C5A9E6A28836FA7422745B3D41E7CFBB7C5C3E47942012ACFA6FE339912F5415B4DF38D37BCAFA463A803B974B86BC940E4D7F4B2D61`,
		MD:  `B2B28E49E2A94FB1A50BBD68A0FEAAC69F148B4455E4E8F6373862F1`,
	},
	{
		Msg: `29D825E3AB7B023921EE6904FAF5DC525A363A16BC4B85D5DC16E8A1D454F94B85753BD538ACE7268542EC01752F3438E3EB291D2882783AF7E71442858A7F47ED2C452CD4F2C95E48B943BFEEAA0AE35D0E0F7FC5`,
		MD:  `A2F69C0DF475D98C2810EDC7343ACEE639C345E1C21094D9F8B276A2`,
	},
	{
		Msg: `63AA1ED4AF9A16DA1447BB2BA549D2C5782E4932FC7825B08DE6BB9046F523DC17A5EA3C1780071A1BD910058752AE5364A14F0BB3931A02015028433CE2F1B4C5F104F7C6D0FC8FC6BEAF447B65A8C8DFA7D89A869A`,
		MD:  `91D989FE0120995EF276B64E50CE310CBCCAAE0FB940E768B471E0F6`,
	},
	{
		Msg: `882B2ED1924098DED1979614D6CD6692DC7F3E7DFEB0288CA86A10CA4CA81B9756743BF0614D01F1ABCBB298C6869CCC9133AA9E671037CE9683209113FD1D916C255136FD331F24679940053EE9128802F333A9512D4F`,
		MD:  `A1E221130700A59AD8550D1A10EE81ADDE778E9555A6F2A3DACE49F1`,
	},
	{
		Msg: `345A59368495613291B286972ACFB99B7EED4968E1A5BF70CD1707012724D019C4DD347B2C3A9C029E10DE36543DA91D07252354B663D05DA4B8B63BA047BF3C88055E17106B7536C06E77553CB0CA4CC8F86F290BE0E482`,
		MD:  `55F73F706152CDDE7D5F27032F2FFB84BA7FE78C85F024F2A5BE43CC`,
	},
	{
		Msg: `13CD490A9A24C36188BB1E31B98EB676FFB4971D6B91E3D1C84C74126983C387253E08959F22A330CBE8BBF1AD0D256964CA87B66846BD53100B7B3557B63413DFF33B0536A9295059727E80764E8BB8F6988F47BB39207BA6`,
		MD:  `79FE36E2BCE0932876AC5324508E5C81F6F53D908DF536F3FCE6DC87`,
	},
	{
		Msg: `1BCCBA5DE101CE4BD52A81E8211B6A3A43AF069CBBAF0BBE06F6F8EDD0D3F19466AFCA036AF06270CA1D323BF24225774196AAA5BA5B08C2FD06F38575A05CE2D26C8DD2BACA7456BB5B3FB985662A16567A951BB5555A0C6305`,
		MD:  `62B7AE10EC40F096D62A72361F220DD7CA1264EF1C167A303CD492FB`,
	},
	{
		Msg: `5D0B8F4D6596494A2A513956AF06912DD1EA404F202269A5EF851366A7A4223C610786A3B3BAB292C9094F5EB88732DB10558CC4B0A0615DDF6C3C56989020F32402C16AEF55F48D948065A83AE188AD2382D28E829A28EB8D5798`,
		MD:  `EF511143A2E773351AA5522DB276EA3D34A6E3298C91076A7657F327`,
	},
	{
		Msg: `EDD4DA4772D184AA6AC070F980B140E1FEDB1D7AEA33F978753D17A59722CE8FF3048A7D008E13A86ACF12A8B648D12D264342EFDC6743E8418A20C0A8C659892FCFA6ACE683392334C343E97682B9AFE81B58D1C9E8E788BA10BC61`,
		MD:  `C0FEB03F5A75A1ADC1CA413440A28FA614DB3B2FF74D080DAEDAD381`,
	},
	{
		Msg: `68A081EE6ADE38BB6DDA7D838E450BA4078682D10E0B6807FADD6DC2A878F3701DD9CD266CD84792DAB8673141FD73C1B27DB3FA5EBA7879B76F5EFDD269CDD51098C699FD66DC197847809315928229BF8FD591A2D1930F48C640A001`,
		MD:  `8C8CFF8C6409921C9416EA0436AF1CFACA722A01C5A76194122AB87D`,
	},
	{
		Msg: `8C53788906F2AFE530FE388EAF4F12D1AD41C276E36F00A0F7BEB1F4099E15FD898C42CF4BDBD1EBDDB5D763185D173FF5B6751EA209AE21D503470461FDC38E0781192C8E0965D91509DE40A930082513EB2B01F32C03EC1369B9C0F5E7`,
		MD:  `B5AA2E1D1CC6070488F48DFBBE430C71392682406DCFACDEF7460273`,
	},
	{
		Msg: `127CAE38E7E6F703D5758203356D8A524E4CAEE8CA65FF5E67CB58F8EF70EEE36A8E27BB9246B0518788CAAAC432E44B527C6FB33B93EE2A42844CB9D15C737A338D1A120394DE7BE5D9950774D7DFA17B7CFD5667D3E715FCE098B348CCAD`,
		MD:  `508E841CB6CA865D07F22F4ECA5E1A059D418C8D3868CEDFEACC1295`,
	},
	{
		Msg: `32D5422E432318AC0955328C202BA3A2DE301602D9356D12A9FDA7E4EF3A07C396B8720C58D91A4084415956FFD2DF5689DB1023428CA8986DAEA1D7739156DE2B45E85C580CE1C27410179CAC52D4222BE3330A9B56CDA1E975951A49330126`,
		MD:  `F76A2E5E221543A23064C1E0B3DB9A8E5B92B6C2E93D9C73196714E9`,
	},
	{
		Msg: `7BE26BC54CA6B65B5F2A34062D91C72F4BB0475FEC7A2454A6BB4D55790A50C20DF8C2EF3B8EE132CFF7783B8C7903A471AA14C93BE0F5DA2BE98F2F09FC8ADB69F0B6832DACAA714BDFABE4D1F1CFAA263961F19270885B0065245003F3F99818`,
		MD:  `1958EA83D0F54316C3DB31CBD8E5DF43B726B7CB04C8DBDDDA9EC1C8`,
	},
	{
		Msg: `300F292BE8583D5C4A94A802E07269F0C5F87FF81224FEAE99D4616ED6E322A12295883B2A35C188A0FC30293035085AB69F1A76AE13EBBCEB660A538D91A08C2D0F49685A87225A49AB244DC860B0E9741F6A122C9B4BB2C5A4F2B9FF4B3B62B847`,
		MD:  `041EC3AC5C2EEE6AA007B0AB84AAE3ACEF750F567F5835AA7E4D8BB3`,
	},
	{
		Msg: `13B760978B57AD692B77903B8827674018D6740546198E54C2A816591BA92D038F3224AD220934FEEDAE700DB99BF7C9CCCA6E990378AE3F2129B90267ECF30C65E08F110FF7473113FA50206ACD8D842DF5039EAB578D2DB842D08F48B0722A3CC987`,
		MD:  `23866599F0C0FCEEE7DAB168EA92F36AD337E9E792F4850757EF3950`,
	},
	{
		Msg: `0A7BA94CCB456318BAC4468F37E45F195E653E93E62FFD4448E57F1224BCF6F5E18EA8AD952AF2B024EF8D67AB9910991F5AB628312D62FB7567E03BC48E6EC2DF6F7BD0546BAE00A6122A7BA8F3116C2AB146B4D75A78DE5EF55C55D05A39D44947AA1F`,
		MD:  `D35D40D4F8E10A6ED977B6D90063F5BB5C489AD8F9FF0F5BA9C648BE`,
	},
	{
		Msg: `D3AFA0F18755A635F69A6978C793779FA15F47D54BB87AC811C760156A3DBE00D99FCE6911E188FC5C3720A73D28D05154E4D7C983D7C2FC4F7EEA6C897A3EDD2DE342518ADE5837A3723FE0D697D932BB9031A04F7BBC5A5227DFC0ED25B21D982204923B`,
		MD:  `F33060EB903547370464C50EEB2F61CF23FCE6CE575F8454DAC8DBFA`,
	},
	{
		Msg: `749F5093F2CC0AD8091AAB9D0C5C44F7CCF148E16948FBDF67A228E2041A126CC9E3AEF460D20F8CA51C22C1FF609BC3F1E82FC6EC2C04A7689766D3338EBB6508AFB98E6356D63A946600F012468FFA132978C75993C38D4F193F325468746D5788D929A462`,
		MD:  `6240FA30FE303718BE231D44FC755212535936BAF4C9BE51D72E6A0C`,
	},
	{
		Msg: `472AFD80AEF2ED1507E6593D89AD588520D945AE1BC912AE24CDDEF8768CCAFD4AA309F3928F37D5EACDECABCC3F900B03D271E0625E126AB5CEC786F8807E30C7EFE046168F998E22F8F3FBA747AA2C54C1FC367CBAFFA30EB637BBC7CD6802A654A901A1B9BA`,
		MD:  `BD3C245AFA1D6D375A5F4253BB951380851DBE91FB2B6A0B7B5CC2C9`,
	},
	{
		Msg: `4996AF97B6848B4DF1FBEE0D38D1A57F6217CD581D4B3B2F7BCF1B8DAD9AD6430E2E3A0063CAD52260E0A1CD6FC9E73AAEB85D5236ED917CA32B8BB3D16BAEDBE3D2E0558AE1EBD184EC226362CD47DA8FD9AC4D8E1399F683A0EE43F57B092BD26AEABFB3319501`,
		MD:  `932DDC4B37E74AEB59DFD50BB0F6F92EFE4CA42DD19989D0A386D631`,
	},
	{
		Msg: `6205B8F89A438E040D620370AB27194622178A67A132AFDAE2DEA06567229D5A0414ADB49D958462313BAAE3178FDA138879A9455EF868911B164A702D1CCA86710247BEF529F19C024DE9CF5AB5A406BEAC9105B335FBF886FB30078E87E32AB1B5CF3F29699D5A88`,
		MD:  `31FBB8DC7D4F1B21444CE7D62A597225C4A312A30FD3F61E827335C1`,
	},
	{
		Msg: `7E085CFAFF3FC7CB9A8AD0295131B62A739A5311CBB1F1CD178C5EC2BC5043CCB74EB181D772FDB729DECE47B4D2741067CF7B3E50FC80CC0FC8444789D1FFD5E105F777A406508B43D32695C7A6D93EC018A89D7E8B82B684C918403AE937ACE2735562CDC46A4ECFB3`,
		MD:  `BEC73C4BA94904835BDEE2EC8A981289B5C775C336B85D5A9F8F191E`,
	},
	{
		Msg: `019BF0748EBE64BD02B7BE289D9360D9033636609357A4E0F4CCBC71F008B5C486E298063A983D06B0137699796C8A3E2DB0E5FEBAE00D96FB3CD8D88DC1DD4FDEF1282286D2323ADB3910A7BF48EA40BD36AF8A554E07DC933B89C0E7822F9B305006911A762AA3AAB2F2`,
		MD:  `800C6E01DD9F865B4BE3621E983D93A4262C3CA92C5A299F14FE9F46`,
	},
	{
		Msg: `1A17306829DC50FCCE0BAA79F9E45037F0A65293F2A5889C1A9D8B0F87460586BA282BEC3FB187349085FB4A3615034A31733AF19A65205187B81AC82FD6BA324F9B415ECC6562BF996A4964DB4FDC2A4BF64131DC97A08C9F122A3DE3835A60A0177AA6161F6A05B1639DAE`,
		MD:  `53E98A86C37323C80834016AE74AAD1DBF55027CDDBDA6E1B05AEA1B`,
	},
	{
		Msg: `30B9644006D66F3B006390B72954C356EB756C777645DFAC6608256F9F18FE107C7C22F22ABBCCF74EF27939EE25EDFB8ED378DF99CD775903481D632D30747209BF9EA97E4F614D0E6CE37AB62EB11415C506FFD1564B933706DA4A215407BC1DAB61CECB37352687F954E248`,
		MD:  `D40C010B52B7883AA0A2079AFBD5BF78DD7C6D4625577FEC3DFB9140`,
	},
	{
		Msg: `161689C8759F41E6886143486AE28D442BC0D60168F2C7D91E506802BF56D96E862542EF08CA0C9F0DC3527F9DBF925E5EA6A9FD1DD9FE70AABBFDC7E970C7EE26E527D2A910C9C58A79A624C7F9C33794C11B58F29FD958837DCB1E9D2B14919D5E21B97EF8DE20AA89297F0F16`,
		MD:  `76F5B4F69F54A5876AB7A812622D628313AFBE3C1E723E9FD8E6D2B6`,
	},
	{
		Msg: `8796761A155C10CBC258529EBC4F030D281F949816EBAF0AB47812FDC1D79EA04545ABB8804F9FE5203ABF27B7A145A06B60B29CF0B165CEB3956DDC58769C1F670ADCA4427AC6765BC4711FA35E886C461CA0719553C8CD433F612FD5A2301954AEB1E44525AE05DBBF118E6ACCF7`,
		MD:  `E1173DF1D8969293C4362AFEE0E1FE71CFF6A741DA23303A9815DD66`,
	},
	{
		Msg: `F81F1A051FEF19FF8ECD3D96FD20A4821874F9D1826F987F1A07F2467DF333B7CC5F238AAFE44A3399B758B64E21C4BBA32FB2E585AA1A603D4927B9F226DDF3706FFFD7285E308B776B58662BC0490343EDA394F6DDF479F7082412F8EB11D16DB5769392441130AF4F05509AE57D81`,
		MD:  `91A2229F23197BFA0E0B43FBD1B72D8D165FD9835424A443B76C5966`,
	},
	{
		Msg: `7E169287DCF06FF3CFADB2F8942A1679499B4F66F57E3649C9CD261CA678A65EDF1248663FCAD1F64E9E9913D0CB35148F8854FA5AF3E7FF4E4C857254610DC1EB498DD8285562A9FA44A10CB3B6B6CAC4E44BF4C88D6C7AC534E8F8AE88FAD32C2428859D786C355E07928B45D4E1C272`,
		MD:  `10A976B14E7361553A6EED6324C0F1DA03D11A6D0ED45A04127CF30A`,
	},
	{
		Msg: `16FDFF48330CE1B6EC1859BBD5B4BC2909E316A91B35C2FF5033378F0BFD25208832C1E1CD105CA9DEDD3B833185A36A6184335126067767854DB99844541A57A1B2E29B44C0152771FE470F82B3B07895D4839F82277C1D10A8B869682CCA9EC9440B7E7E6431B70D452466B2870DDA0B9B`,
		MD:  `BF5E885C0328F3024EF1F0364F7BFB1584050019DFD86D7DE4FE6623`,
	},
	{
		Msg: `AF9A15E0FAA9570CDD085DCE5D7D68B2E4936444DBE21A787A90A749DF182DA15E4AC8791994C441D07F0D82FA7CE6C214E1DEA8392F0060F3964B4151BB62AD39B884390064BD0E38D65E90C924AC3D4470D5D7732628C5357F84B7884A1189E8EB2F6AD8A4AC3DBE133A3FF3FC20EA1B492F`,
		MD:  `83CFEA1A360F13176D0AB2826E9FA6FEAE751DDEC4911B2FF018F259`,
	},
	{
		Msg: `1AB24FF86FE8533CF2AF4DA4FB12EBB252C8D9F894094088D8A08DAB934273D6B91E33D925B51B5A3BEC4C82CC5FC1F639F98D09DCC26430A3F75DB6FE60F6B0FCDE5934B144DEC6DA67D013C47974E269924747DC908819E821F6A7F41EDBA72A1B496FD93CA468174A60BFAA050C5D29BC901A`,
		MD:  `2A5A1E24CC34E6C21419B2E1FDB847A9E43629D58EDBA54CA7804108`,
	},
	{
		Msg: `5E1CB32DEF46E0E73D2429A28DFC8CF46645752839BA85EE48AE9251C09FEFEA280DED53FF8CA1DFA119E8886BAAA9291F095AEE8AE1C60864190DA57762EC1F040C3E73621EC3F6E8A521430D4C2878CE2D5AA250E7AB75C7F18541962E0F21B8C1AF765E22CE93D971DBBABB0B9248395C7E2CD4`,
		MD:  `A6E91A064C3C6F884B82655CEA49B836AD56C521B89DC979ED61CB1D`,
	},
	{
		Msg: `F4BC9D39DC030D652AE1580CD347423CC148203E251ED5BD0F6E00293CCD745E2E82E981828DF12B7FE92DCAEB0FC04422BB04DD69F27C5A7F440A66754F6236C084BF698D2ED53254AF06B8B00DD8BD6EAEA77EE36097DB839DCDB8B00CDFA3D04A88B7A2EBA1C0D2BC4CA0CB4B73F1CCAA71EB3E6F`,
		MD:  `D8D5F5BC82159CB99CE6D209A612DA36B76A766EA7197F3419F5D058`,
	},
	{
		Msg: `0E8845FFB1E7950DD1EE875DAD346BCE3D1E20976848E3D177266536CAEB2028C9A9EAF51EBACB5922600C4494A890DE919468C8C078A1A5CEA560DA1A453291B1869D047CFBFBF8B83F8887A7DBB9238569CFD8BEFAADD15EB6E61EC9C21BB5AC6784B63370CCF37DD322AF4D948DA602E77F4EC995AE`,
		MD:  `F7EDF07A8C81E018FE8A5FFE1BC4CAE5E8533D755D32B2CC48A3F873`,
	},
	{
		Msg: `5D11A01F8B3F2A8399A393C792ABE2E75157660343ED9DB46EEE54A0B248104C4FCDAD32A8E44514477F538F204D18F31A303E52CBF1BCD78333F9598DA53430B17C2348FD21DA74C626EF04E97FDC9956EFA438D0C0BA9EEEB71ED36647C9028A054B88135CB34B1E82AF4D74AAF4810E7CE0DDD9C25EB0`,
		MD:  `925823134BC858DAC51D8A7C9F675FB73949340DA03374BA0925325E`,
	},
	{
		Msg: `B6640B36C5969769E3506D1B71FE9642587ADA93493729524ABF3A686E5F72B6E32007A932FA650EC345D128A0704C007AF5BAFDF18C43B410A0FA0117FFBC577ED55B5436DBFC9932CCBED826199D2E8FDD95B4F3130A516A0767D25D0A063C8E6A506058A125B139CD7E538D6050749C19B4DFD73346CE5F`,
		MD:  `7987EBF339CAF5403602572E8DCC2DC54157AFAD8C66B2A947E617AB`,
	},
	{
		Msg: `DCCA00C7E8160C1BC72CF521B4A6FEDE3057C05CAF2F79AB915FA72653BED6F7076BB6DE01B4702ED4373A187D472A2BDD92471C6D00E6C32DDE3467DD3B0F3B70921A28BEAA10CD79BB58C562D2DA28BC761254018C62EDD578AD2238852132FF87542F85BF93039EFE77B0CC0DB8DCB64D468CD59B96307176`,
		MD:  `2CFFE2F2050839A7FACDC68D1E9996B30AEA55FC0E773C0B972A4349`,
	},
	{
		Msg: `3AE0CC64B5E59D918D7BBFBAA2D52B49BD6224C26C4FB15D0E4D50F52A94E920052701A0735A30FAED92D564E44D08E45562D0B79E8996F0901DDA61C52125BE33823F00029BA1089BFE2D8502721D1765383AB5F5934E4CF487DBA80FE00D98EF76EE4D0BB0D904C04731E5A34C522A166DC8831974C0B96310BB`,
		MD:  `E451D272431A49077F2D07FBDEF2D0C9D8AB3BBE4672C16F96E6CB38`,
	},
	{
		Msg: `92B2D88AF65F8D764F359FBF3683F3484B70E490648F88E4EBB334D3810542E3E75311821FAB4E3FCE29723AE50D8CD9801041BD0FF911F5B36BDC71EB01A4B96F9D94DF5CFC41B1928E31345A859A4E4D2C6BC3717B9063A4E872813E7FBEA6F2356887F57AEF35620DE294DE14A30E786C6B88D23E8E9586D301EA`,
		MD:  `8633C83DFEEAF5D637096F28A58236A349F8E75A1C85B62BE9388DC3`,
	},
	{
		Msg: `C3A9191337C69564324CDF31084E8F42B7E23181C97748F67C6993F521D38E96C2A53C6607F6245C5A9B990F4AF83F252D16CB0F0AF60AF4B62525EA3C91236C829BAF092A6D1C59B78EA635F1B3FB278AD9FA200190BA69C79ACEDCA1B4C7544DB4FB6268D485428F4E8D0DC529CA2F1452441C77EE5F3FF8EF153316`,
		MD:  `3D62B10E26D8925105E0C96887D8B6D7BA198ECE0C68B31352F8E12F`,
	},
	{
		Msg: `E583366917A7DB89D368FF2B02D6FD8DFFB42615AA346B35B8015D914D825864CC5EFA1FEE3BDA5411CC77D3D30A6F13E0C1F2D3084F750095A8ABDADB9EB153106089B75E40C99237B994CFE9B0CA1D72FA7C97951FA3A32D5D111AC2F0A7AA5A91C979913D0F0BEC373D666B5076745D66ACA8782B5BE42927779B2901`,
		MD:  `28DE8B8918DAC38349350A6669C32D4B610EED6D900AA070A63286B8`,
	},
	{
		Msg: `A33351FAA3ADAFBA1487D94F5AE18173DBE04A0D6B50DAFB068D310FC5A339570BB49EA2FF9ADF3F40CEA608D48855A06128E1F89FD5E3180D64EC1CE37B197C48F6E891C6894A55FBA7FA6B91A30C5F7016C92E3180BA09D6320356B96BF3F8C3C1D8D0DEC21883ADDC7D7E7A7AD1387070ABB70A998856CBBE1DC530BAEC`,
		MD:  `C7CF1AEB9B0FADDB75D50EBA81FDA7F7687C93656430D4615270FDE3`,
	},
	{
		Msg: `51E3DC1666318D7052BE33D84B7755CAB4C9FF2E9C1068748853D2E3722EB22E110F11495BA6327A499675B3DA717B1F5AE38BFBE42D34A726D534FD649C99E1BFB12619B824F860F09D6AB5314E35F9B226DCB997484749B8BBAD12A363B82CE430E43084952947CF3CD9875AF8CB8045FBB911BE49E781C66C05CA38695738`,
		MD:  `5A409D2C7B9CBC4E45DB8E468184FBC57239706F9B9264C0B3DCB399`,
	},
	////////////////////////////////////////////////////////////////////////////////////////////////////
	// 암호알고리즘 검증기준 V3.0
	// 테스트 벡터
	// LSH(256-224)LongMsg.txt
	// 용량 문제로 일부만 추가
	{
		Msg: `CB8CDE5AFD7AC4D740CCC149EA0D80C4910DE2B1CD92E91F3E4BD55AD1058CBF40213FCED8464B817F535DDCF598DAC65A457E17F12243F223818455C77262EEA035673127F569A78BB3B189F7805D8D13EF0300CDCA971559048F7D76EF4437F4AF9CAA73EFE6F16DD20576772CD1C2C3F3BA93EAF0AF2FE8BA3EF968E5FC4243772761439B8A062C9AA1CBFD06FE0C7339EBF9CFFE51E9D98DD8F225BD50B69575507F1F631C8FD0F1CCB11175AB142AA8DE59045A46EB3364A78F33DD093CF2915E2C8FAE663362C1D0503CE1A282F02ADADC126E55DCFE26F1F91CAFEF7B63B39B`,
		MD:  `7617289F423459CD2DEA4E69AD71BA0EE17022B75F9F1532D3C4AAF5`,
	},
	{
		Msg: `8F1CE52E9BE9F1F3CF05B2A9FDCDA629AA80A0486543BDFF576899C91CEFC34DD14D6E3D6E6D697E56F551892DCA0511EAD65AE62D091019D29D0561C69FA9BD19ABB25B55F711B90F92240ABBEE2DDBC59D03F6DD58078692455A7D227A4D3390BD194FF88DD1EB206794F34F4207CE625C02411C388A90BEE700E6D78897F65E0B0CE3FE3A915DB1F90A0E90CC7B34EA9BC18F93B2817D7E0C5F648ED17050AA1DF1DD10033958EBD2EC212796725384E2A7A7E9CDD496FC3BDFBEEF5EBE879C7C3106D5F2B122F579A3F5BBA8D27559BA1C53C7926C235DFCE8BEA1366BE45EAF3526F60EE103F8769852F50885E090AA8859D5072F6CCCD7E22A4D615DAF153E64051A5FB1451C5232FF7DC071DC523B5482BA3606B86F5535CB9AE87E30ECB2256AEAEE0A2F8894D8C6FAD880B3C7F3E6961F27D9516EFD496931D2B4AF0991E21F0DC2`,
		MD:  `811EFB397CC243D6CFF9EB87EA4F206889A698386F71960070C35CC1`,
	},
	{
		Msg: `D20865C4F46D155699AB165DA75DABE1917DF25886CCB927E973946602EA2CE3F319DA6AEDBD7544A30E68FEFF9F076C148523EC55BBDBEF04C5B6B6ED94EF5A53AA100F2B7EC2A98822E4605723DAD04CB128441CFF3F33EDED960B0AE5FD2EDC022F3A1678A68E2F30B30A15B2CD1520482AE544608CFB8C335EECC3A4FF67746B5E741672C7FA807FFC84A31488447794505936A66B4DC9DE33E27F9B9C0E052BC5449336CDF56359E8576711B2653BDBC327599983338C383F74A7907D2A768B8D32F58A6188C0049F0ACB70F37E5366AAD714017DB4BD87AA2AA14C4AC3AFD3DEC6A4372BB97FC9482535FBF599A77D2EF1D1A500D844149C8CD697AAE2984BDF890805D29288F00B300F785DBD1F6777FEF74EB5A809273C23AE38458D193BB90289BBFF1BC3C111B4C0B0D6F3A36EAC86EEC4432BF408B27692F9C4CF1AEA93BA9023595B19848239B06B06421BD8F50F1DCF546AEEFE270DE8A0CFAD83A196378303895A70808645477195B26FEE7B23EE378D6BDD52C37119F60D313CA8EA03CC243620B2FE4362ED892C4B87F86649C8C39838AB4DAD298EC327622D47B6A6D14E09B6C5`,
		MD:  `B7547C987224595BD6057DD6CC2C6799140B7FE6E33547D857A1BBB3`,
	},
	{
		Msg: `46E4170A7D73164B3EDE0A133D1DD93E340EBDAECFC4483FC622A6FC48A92393A08FEC0613111AA4070AEC386BC35480520DB6E0E18DEC586C588DB4DCC06F02546D6A4914AF5F789D12799D1533B357EED39B8EC32A5561827DE260D0F3A3FCA9A5183B47F4EEB1461F39AB0515FEAD4F87F41F3B428B4350BC2E6AEDECC568E85F7F457689CE354B59B3CA02103B69BE524AFB127D97E41E3BF839FB53E02D5AC525CD671507EEF4E76F63D3CC7473835C238A1085802E3422C835C2CFF93544FC943DE341A2C288F1F5DC41F2B0B3E6CC0A33FCFF4E08D99A26C609091A67F12F54FBB2B5A69B51A1CE9F142AF8112FCA84619F960A5957CB9B549AA94AABA884EB719C181B60951D8112141B5375A5801B79C2F0BB0BF4DCAE473B5692E9B023E30569120AF99D8D969F086EC9C7911456E42BB56905F9F5E807B5B8F8095235C421E04D7A4DB11A95ADBACA62EF3BAEBD0AA38E1C2FAE40D0FCD07887F4D6E2142CCB6E734618EC07A3F0D827D5EB77D953F123DD715AE2EE8F533E449084515D8EF020FDCB1C2A74EB73401F61E9963128DF1BB3B34606CA2608B139C6A3AB26AE170921C503FC62EC0BA9527C7C344EEF331EA7DDBBD2ED2BB57A6D7F7D17F8F60AD1E61A168B5B0E7FBBC90CEE79B612B6D6C0D7FF6EDE042341E8A158BE5ACD902155B39DFFE6B9991F8BFA858CF3F730E806895A03251B1AADB3157DC8D49B70A17478E0808C55`,
		MD:  `679D15A8D19C91F5806BA63B3D79BC0CCC9CF24DE13822A82888F325`,
	},
}

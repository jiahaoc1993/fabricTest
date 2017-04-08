$(document).ready(function(){
	//getBlockInfo("11")
	$("#24BI").click(function(){
			$(".active").attr("class","");
			$(this).attr("class","active");
			getBlockInfo(24);
		});

	$("#23BI").click(function(){
			$(".active").attr("class","");
			$(this).attr("class","active");
			getBlockInfo(23);
		});

	$("#22BI").click(function(){
			$(".active").attr("class","");
			$(this).attr("class","active");
			getBlockInfo(22);
		});

	$("#21BI").click(function(){
			$(".active").attr("class","");
			$(this).attr("class","active");
			getBlockInfo(21);
		});

	$("#20BI").click(function(){
			$(".active").attr("class","");
			$(this).attr("class","active");
			getBlockInfo(20);
		});

	$("#19BI").click(function(){
			$(".active").attr("class","");
			$(this).attr("class","active");
			getBlockInfo(19);
		});


	$("#18BI").click(function(){
			$(".active").attr("class","");
			$(this).attr("class","active");
			getBlockInfo(18);
		});

});



function getBlockInfo(num) {
	var url = "http://54.179.182.63:7050/chain/blocks/" + num;
	$.get(url, function(result){
		$("#stateHash").text(result.stateHash);
		$("#PbHash").text(result.previousBlockHash);
		$("#txid").text(result.transactions[0].txid);
		$("#nonce").text(result.transactions[0].nonce);
		
	});

}

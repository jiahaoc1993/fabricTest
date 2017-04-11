var url = "http://54.179.182.63:7050/";

$(document).ready(function(){
	//getBlockInfo("11")
	var height = getCurrentHeight() - 1;
	init(height);
		
	for(var i = 1;i <  height; i++){
		//var  next = i-1;
		$("#init").after("<li><span id=\"block"+ i +"\">"+ i +"</span></li>");
		(function(){
			var tmp = i;	
			$("#block"+tmp).click(function(){
				$("#blockNum").val(tmp);
				$(".active").attr("class","");
				$(this).attr("class","active");
				getBlockInfo(tmp);
			});
		})();
	};

});

function init(height){
	$(".active").text(height);
	$("#blockNum").val(height);
	getBlockInfo(height);	
	
	//alert($("#block"+height).text())
	$("#height").click(function(){
		$("#blockNum").val(height);
		$(".active").attr("class","");
		$(this).attr("class","active");
		getBlockInfo(height);
	});

}

function getBlockInfo(num) {
	var u = url + "chain/blocks/" + num;
	$.get(u, function(result){
		$("#stateHash").val(result.stateHash);
		$("#PbHash").val(result.previousBlockHash);
		$.each(result.transactions, function(index, value){
			$("#scroll").empty();
			var html = "<dl><dt>Transaction" + index +"</dt></br><dd>transaction id: " + value.txid + "</br>nonce: " + value.nonce + "</dd>";
			$("#scroll").append(html);
			//$("#txid").val(result.transactions[0].txid);
			//$("#nonce").val(result.transactions[0].nonce);
		});
	});

}

function getCurrentHeight(){
	var u = url + "chain" ;
	var h ;
	$.ajax({
		url: u,
		type: "GET",
		async: false,
		success: function(result){
			h = result.height;
		}
	});
	
	return h;

}


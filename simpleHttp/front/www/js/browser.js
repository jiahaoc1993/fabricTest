function getBlocksTotal(){
	$.get("http://54.179.182.63:7050/chain", function(result){
		alert(result.high)	
	});

}

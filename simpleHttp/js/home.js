$(document).ready(function(){
	var user = getUrlParam('name')
	$("#accountName").text(user)
		$.post("http://127.0.0.1:8000/query", {"user":"Luke"},
			function(result, status){
					alert(result.usd+" "+result.gpcoin+" "+status);
		}
	);
	
	$("#topup").click(function(){
		$("#txType").val("top up");
		$("#amount").attr("placeholder","100 USD");
		$("#fee").val("0 USD");
	});

	$("#invest").click(function(){
		$("#txType").val("invest");
		$("#amount").attr("placeholder","100 GP Coins");
		$("#fee").val("5% USD");
	});

	$("#cashout").click(function(){
		$("#txType").val("cash out");
		$("#amount").attr("placeholder","100 GP Coins");
		$("#fee").val("5% USD");
	});

	$("#transfer").click(function(){
		$(".dialog2").show()
	});	
	
})


 function getUrlParam(name) {
            var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)"); 
            var r = window.location.search.substr(1).match(reg);  
            if (r != null) return unescape(r[2]); return null; 
        }

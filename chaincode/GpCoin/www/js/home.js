$(document).ready(function(){
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
		$("#txType").val("transfer");
		$("#amount").attr("placeholder","100 GP Coins");
		$("#fee").val("0 USD");
	});
})

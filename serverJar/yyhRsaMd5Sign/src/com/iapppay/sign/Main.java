package com.iapppay.sign;

public class Main { 
	
	public static void main(String[] args) throws Exception {	    		
		if (args.length == 0)
		{	
			System.out.print("");
			return;
		}
		   		
    	String content = args[0];
    	String inputSign = args[1];
    	String publicKey =  args[2];
    	
    	//String content = "{\"appid\":\"5001324819\",\"appuserid\":\"YYH5350724\",\"cporderid\":\"2016062316374611198\",\"cpprivate\":\"YYH5350724\",\"currency\":\"RMB\",\"feetype\":0,\"money\":1.00,\"paytype\":401,\"result\":0,\"transid\":\"32011606231637464769\",\"transtime\":\"2016-06-23 16:37:56\",\"transtype\":0,\"waresid\":1}";
    	//String inputSign = "RBGZEL518F32D0Htape8oudmK0tWXFfoFrvIMr9No9hoW5QWrP+ODaQ2wK8Hnm7ZYTzO3HyHkm11ZJGDp0aEqExdtxMrnrHh1UPpoRymtWRxfF4xkSfKgPzDTS+unbQvwohgqWK1HpZa4ifPTdUl5Cc0mRBs+O//Xo9LCp/l90o=";
    	//String publicKey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCTCxn4EQusS8pPuJ0Kd6M+Q7elbYCI38pLkgxJkD7PdZevRc8AC9dADzcbbJALtm2InRG+nOyxQmNHF9Nogp0i01QtEvCYyRa1y3EL619w1OWX6wRwY/PIouIOy5zb2X4oLldPTlSY9nh6W7CksQhgMpn2ey1mjY7/6oZhrPtg5wIDAQAB";
	
    	boolean isSuccess = false;
    	try 
		{
    		isSuccess = SignHelper.verify(content, inputSign, publicKey);
	    	if (isSuccess) 
	    	{
	    		System.out.print("0");
	    	} else {
	    		System.out.print( "1");
	    	}
		}
    	catch(Exception e)
		{
    		System.out.print( "1");
    		return;
		}
	    
	  }

}
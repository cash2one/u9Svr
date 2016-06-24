package com.iapppay.sign;

public class Main { 
	
	public static void main(String[] args) throws Exception {	    		
		if (args.length == 0)
		{	
			System.out.print("");
			return;
		}
		   		
    	String content = args[0];
    	String privateKey =  args[1];

    	try 
		{
    		sign = SignHelper.verify(content, privateKey);
	    	System.out.print(sign);

		}
    	catch(Exception e)
		{
    		System.out.print("");
    		return;
		}
	    
	  }

}
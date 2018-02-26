package it.infn.ct;

import java.lang.Thread;

import org.json.simple.JSONObject;

public class workflow 
{

	public static String vmState = "inactive";
	public static String configEGI = "/home/configEGI.json";
	public static String[] vmFeatures = null;

	public static void main(String[] args) {
		// Get input from EGI
		JSONObject egiInput = Def.getJson(configEGI);

		// Instantiate a new VM on EGI based on the input from configEGI.json
		String vmId = new String (InstantiateVM.instantiateVM(egiInput));
		System.out.println("Waiting for the Virtual Machine to be active...");

		int i =0;
		// Describe the VM state - Once this part is successfull the VM is operational
		while (!vmState.equals("active") && i<10) 
		{
			try 
			{
				Thread.sleep(10000);
				vmFeatures = DescribeVM.describe(vmId, egiInput);
				vmState = vmFeatures[1];
				i+=1;
				Thread.sleep(10000);
				System.out.println("Try number: "+i);
				System.out.println("Virtual machine state: " + vmFeatures[1]);
			} catch(InterruptedException ex) {Thread.currentThread().interrupt();}
		}

		//Write the public IP and vmState on the JSON file
		Def.WriteJson(vmId,vmFeatures[0],vmFeatures[1],configEGI);

	}
}
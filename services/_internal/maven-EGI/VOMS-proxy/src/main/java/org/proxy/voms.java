package org.proxy;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.security.InvalidKeyException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.SignatureException;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.List;

import org.bouncycastle.asn1.x509.AttributeCertificate;

import org.italiangrid.voms.VOMSAttribute;
import org.italiangrid.voms.VOMSValidators;
import org.italiangrid.voms.VOMSGenericAttribute;
import org.italiangrid.voms.ac.VOMSACValidator;
import org.italiangrid.voms.credential.UserCredentials;
import org.italiangrid.voms.util.CertificateValidatorBuilder;
import org.italiangrid.voms.util.CredentialsUtils;
import org.italiangrid.voms.request.VOMSACService;
import org.italiangrid.voms.request.impl.DefaultVOMSACRequest;
import org.italiangrid.voms.request.impl.DefaultVOMSACService;

import eu.emi.security.authn.x509.impl.PEMCredential;
import eu.emi.security.authn.x509.proxy.ProxyCertificate;
import eu.emi.security.authn.x509.proxy.ProxyCertificateOptions;
import eu.emi.security.authn.x509.proxy.ProxyGenerator;
import eu.emi.security.authn.x509.X509Credential;
import eu.emi.security.authn.x509.X509CertChainValidatorExt;

public class voms {
	static final String keyPassword = "";
	  
	public static void createProxy() {
		
		try {
		PEMCredential c = new PEMCredential(new FileInputStream("/root/.globus/userkey.pem"), new FileInputStream("/root/.globus/usercert.pem"),
			keyPassword.toCharArray());
		
		X509Certificate[] chain = c.getCertificateChain();
		
		VOMSACValidator validator = VOMSValidators.newValidator();
		List<VOMSAttribute> vomsAttrs =  validator.validate(chain);
		  
		if (vomsAttrs.size() > 0) {
			
		    VOMSAttribute va = vomsAttrs.get(0);
		    List<String> fqans = va.getFQANs();
		    	
		    for (String f: fqans)
		    	System.out.println(f);
		
		    List<VOMSGenericAttribute>	gas = va.getGenericAttributes();
		
		    for (VOMSGenericAttribute g: gas)
		    	System.out.println(g);
		}

		X509Credential cred = UserCredentials.loadCredentials(keyPassword.toCharArray());
		X509CertChainValidatorExt validatore = CertificateValidatorBuilder.buildCertificateValidator();
		VOMSACService service = new DefaultVOMSACService.Builder(validatore).build();

		DefaultVOMSACRequest request = new DefaultVOMSACRequest.Builder("fedcloud.egi.eu").lifetime(12*3600).build();
		      
		AttributeCertificate attributeCertificate = service.getVOMSAttributeCertificate(cred, request);
		  
		ProxyCertificateOptions proxyOptions = new ProxyCertificateOptions(cred.getCertificateChain());
		proxyOptions.setAttributeCertificates(new AttributeCertificate[] {attributeCertificate});
			  
		ProxyCertificate proxyCert = ProxyGenerator.generate(proxyOptions, cred.getKey()); 
		  
		OutputStream os = new FileOutputStream("/tmp/x509up_u5040");
		CredentialsUtils.saveProxyCredentials(os, proxyCert.getCredential());		  
		} catch (InvalidKeyException ex) {
			// TODO Auto-generated catch block
			ex.printStackTrace();
		} catch (SignatureException ex) {
			// TODO Auto-generated catch block
			ex.printStackTrace();
		} catch (NoSuchAlgorithmException ex) {
			// TODO Auto-generated catch block
			ex.printStackTrace();
		} catch (KeyStoreException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (CertificateException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (FileNotFoundException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}				    
	}
	  
	public static void main(String[] args) {
		createProxy();
		try{
			Thread.sleep(10000);
		}catch(InterruptedException ex) {Thread.currentThread().interrupt();}
		
		System.out.println("Proxy created");
		System.exit(0);	
	}
}

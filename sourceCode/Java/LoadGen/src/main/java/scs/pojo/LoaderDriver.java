package scs.pojo;
  
import scs.util.loadGen.driver.AbstractJobDriver;

public class LoaderDriver {
	private String loaderName;
	private AbstractJobDriver abstractJobDriver;
	
	public LoaderDriver(String loaderName, AbstractJobDriver instance) {
		 this.loaderName=loaderName;
		 this.abstractJobDriver=instance;
		 
	}
	public String getLoaderName() {
		return loaderName;
	}
	public AbstractJobDriver getAbstractJobDriver() {
		return abstractJobDriver;
	}
	public void setLoaderName(String loaderName) {
		this.loaderName = loaderName;
	}
	public void setAbstractJobDriver(AbstractJobDriver ajd) {
		this.abstractJobDriver = ajd;
	}
}

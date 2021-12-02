package scs.util.rmi;

import java.rmi.Remote;
import java.rmi.RemoteException; 
/**
 * RMI interface class, which is used to control the load generator
 * The functions can be call by remote client code
 * @author Yanan Yang
 * @date 2019-11-11
 * @address TianJin University
 * @version 2.0
 */
public interface LoadInterface extends Remote {
	public float getWindowAvgPerSecLatency(int serviceId,String metric) throws RemoteException; //return the value of Avg99th
	public float getRealPerSecLatency(int serviceId,String metric) throws RemoteException; //return the value of queryTime
	
	public float getTotalQueryCount(int serviceId) throws RemoteException; //return the value of queryTime (95th), unused
	public float getTotalRequestCount(int serviceId) throws RemoteException; //return the value of queryTime (99.9th), unused
	public float getTotalAvgServiceRate(int serviceId) throws RemoteException; //return the value of SR
	public int getRealQueryIntensity(int serviceId) throws RemoteException; //return the value of realQPS
	public int getRealRequestIntensity(int serviceId) throws RemoteException;  //return the value of realRPS
	
	public void execStartHttpLoader(int serviceId, int intensity, int concurrency) throws RemoteException; //start load generator for serviceId
	public void execStopHttpLoader(int serviceId) throws RemoteException; //stop load generator for serviceId
	public int setIntensity(int intensity,int serviceId) throws RemoteException; //change the RPS dynamically
}

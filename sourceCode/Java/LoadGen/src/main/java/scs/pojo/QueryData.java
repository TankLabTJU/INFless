package scs.pojo;

public class QueryData{
	protected String loaderName;
	protected long generateTime;
	protected float queryTime999th;
	protected float queryTime99th; //99th latency per second
	protected float queryTimeAvg;
	protected float queryTime95th;
	protected float queryTime90th;
	protected int avgQps; //window size average QPS
	protected int avgRps; //window size average RPS
	protected float totalAvgServiceRate; //window size average QPS/RPS
	protected int realQps; //average QPS per second
	protected int realRps; //average RPS per second
	
	public QueryData() {
	}
	public long getGenerateTime() {
		return generateTime;
	}
	public void setGenerateTime(long generateTime) {
		this.generateTime = generateTime;
	}
	public float getQueryTimeAvg() {
		return queryTimeAvg;
	}
	public void setQueryTimeAvg(float queryTimeAvg) {
		this.queryTimeAvg = queryTimeAvg;
	}
	public int getAvgQps(){
		return avgQps;
	}
	public void setAvgQps(int qps){
		this.avgQps=qps;
	}
	public int getAvgRps(){
		return avgRps;
	}
	public void setAvgRps(int rps){
		this.avgRps=rps;
	}
	public float getTotalAvgServiceRate() {
		return totalAvgServiceRate;
	}
	public void setTotalAvgServiceRate(float serviceRate) {
		this.totalAvgServiceRate = serviceRate;
	}
	public float getQueryTime90th() {
		return queryTime90th;
	}
	public float getQueryTime95th() {
		return queryTime95th;
	}
	public float getQueryTime99th() {
		return queryTime99th;
	}
	public float getQueryTime999th() {
		return queryTime999th;
	}
	public void setQueryTime90th(float queryTime90th) {
		this.queryTime90th = queryTime90th;
	}
	public void setQueryTime95th(float queryTime95th) {
		this.queryTime95th = queryTime95th;
	}
	public void setQueryTime99th(float queryTime99th) {
		this.queryTime99th = queryTime99th;
	}
	public void setQueryTime999th(float queryTime999th) {
		this.queryTime999th = queryTime999th;
	}
	public int getRealQps() {
		return realQps;
	}
	public int getRealRps() {
		return realRps;
	}
	public void setRealQps(int realQps) {
		this.realQps = realQps;
	}
	public void setRealRps(int realRps) {
		this.realRps = realRps;
	}
	public String getLoaderName() {
		return loaderName;
	}
	public void setLoaderName(String loaderName) {
		this.loaderName = loaderName;
	}
 
}
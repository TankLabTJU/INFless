package scs.pojo;

public class PageQueryData extends QueryData{

	private float windowAvg99thQueryTime; 
	private float windowAvgAvgQueryTime; 
	
	public PageQueryData(QueryData data) {
		this.generateTime = data.getGenerateTime();
		this.queryTime99th = data.getQueryTime99th();
		this.queryTime95th = data.getQueryTime95th();
		this.queryTime90th = data.getQueryTime90th();
		this.queryTime999th = data.getQueryTime999th();
 		this.queryTimeAvg = data.getQueryTimeAvg();
		this.avgQps = data.getAvgQps();
		this.avgRps = data.getAvgRps();
		this.realQps = data.getRealQps();
		this.realRps = data.getRealRps();
		this.totalAvgServiceRate = data.getTotalAvgServiceRate();
	}
	public PageQueryData() { 
	}
	public void setWindowAvg99thQueryTime(float windowAvg99thQueryTime) {
		this.windowAvg99thQueryTime = windowAvg99thQueryTime;
	}
	public void setWindowAvgAvgQueryTime(float windowAvgAvgQueryTime) {
		this.windowAvgAvgQueryTime = windowAvgAvgQueryTime;
	}
	public float getWindowAvg99thQueryTime() {
		return windowAvg99thQueryTime;
	}
	public float getWindowAvgAvgQueryTime() {
		return windowAvgAvgQueryTime;
	}
 
}
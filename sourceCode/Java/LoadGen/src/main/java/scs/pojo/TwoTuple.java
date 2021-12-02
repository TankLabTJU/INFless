package scs.pojo;

import java.io.Serializable;
/** <p>Title: TwoTuple</p>
 * <p>Description: 两个元素的元组，用于在一个方法里返回两种类型的值</p>
 * @version 2012-3-21 上午11:15:03
 * @param <A>
 * @param <B>
 */

public class TwoTuple<A, B> implements Cloneable, Serializable{
  
	private static final long serialVersionUID = 1L;
	
	public A first;
	public B second;

	public TwoTuple(A a, B b) {
		this.first = a;
		this.second = b;
	}
	public TwoTuple() {
	}

	public String toJsonStringBasicDataType() {
		return "[{\"first\":"+first+",\"second\":"+second+"}]";
	}
	
	@Override
	public String toString() {
		return "("+first+","+second+")";
	}
	
	@Override
	public TwoTuple<A, B> clone() throws CloneNotSupportedException {
		return new TwoTuple<A, B>(first, second);
	}
	
}
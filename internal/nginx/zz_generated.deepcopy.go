//go:build !ignore_autogenerated

/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package nginx

import ()

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CommonConfig) DeepCopyInto(out *CommonConfig) {
	*out = *in
	if in.Resolvers != nil {
		in, out := &in.Resolvers, &out.Resolvers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.WorkerProcesses != nil {
		in, out := &in.WorkerProcesses, &out.WorkerProcesses
		*out = new(int)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CommonConfig.
func (in *CommonConfig) DeepCopy() *CommonConfig {
	if in == nil {
		return nil
	}
	out := new(CommonConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReverseProxyConfig) DeepCopyInto(out *ReverseProxyConfig) {
	*out = *in
	if in.Logfmt != nil {
		in, out := &in.Logfmt, &out.Logfmt
		*out = new(string)
		**out = **in
	}
	if in.Rewrite != nil {
		in, out := &in.Rewrite, &out.Rewrite
		*out = make([]RewriteConfig, len(*in))
		copy(*out, *in)
	}
	if in.HideHeaders != nil {
		in, out := &in.HideHeaders, &out.HideHeaders
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ProxyHeaders != nil {
		in, out := &in.ProxyHeaders, &out.ProxyHeaders
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReverseProxyConfig.
func (in *ReverseProxyConfig) DeepCopy() *ReverseProxyConfig {
	if in == nil {
		return nil
	}
	out := new(ReverseProxyConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TCPReverseProxyConfig) DeepCopyInto(out *TCPReverseProxyConfig) {
	*out = *in
	if in.Servers != nil {
		in, out := &in.Servers, &out.Servers
		*out = make([]Upstream, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TCPReverseProxyConfig.
func (in *TCPReverseProxyConfig) DeepCopy() *TCPReverseProxyConfig {
	if in == nil {
		return nil
	}
	out := new(TCPReverseProxyConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Upstream) DeepCopyInto(out *Upstream) {
	*out = *in
	if in.Weight != nil {
		in, out := &in.Weight, &out.Weight
		*out = new(int32)
		**out = **in
	}
	if in.MaxConns != nil {
		in, out := &in.MaxConns, &out.MaxConns
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Upstream.
func (in *Upstream) DeepCopy() *Upstream {
	if in == nil {
		return nil
	}
	out := new(Upstream)
	in.DeepCopyInto(out)
	return out
}

#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/in.h>
#include <linux/ip.h>
#include <linux/tcp.h>

#include <bpf/bpf_helpers.h>
#include <stdlib.h>

struct {
  __uint(type, BPF_MAP_TYPE_ARRAY);
  __uint(max_entries, 1);
  __type(key, __u32);
  __type(value, __u32);
} port_map SEC(".maps");

SEC("xdp")
int drop_packet_by_port(struct xdp_md *ctx) {
  void *data = (void *)(long)ctx->data;
  void *data_end = (void *)(long)ctx->data_end;

  void *size = data + sizeof(struct ethhdr);
  if (size > data_end) {
    return XDP_PASS;
  }

  struct ethhdr *eth = data;
  if (eth->h_proto != __constant_htons(ETH_P_IP)) {
    return XDP_PASS;
  }

  size += sizeof(struct iphdr);
  if (size > data_end) {
    return XDP_PASS;
  }

  struct iphdr *ip = data + sizeof(struct ethhdr);
  if (ip->protocol != IPPROTO_TCP) {
    return XDP_PASS;
  }

  size += sizeof(struct tcphdr);
  if (size > data_end) {
    return XDP_PASS;
  }

  struct tcphdr *tcp = data + sizeof(struct ethhdr) + sizeof(struct iphdr);

  __u32 key = 0;
  __u32 *port = bpf_map_lookup_elem(&port_map, &key);
  if (port && tcp->dest == __constant_htons(*port)) {
    bpf_printk("dropping packet on port %d\n", *port);
    return XDP_DROP;
  }

  return XDP_PASS;
}

char __license[] SEC("license") = "GPL";

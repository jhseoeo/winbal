#include <stdlib.h>
#include "./lib/vpx_encoder.h"
#include "./lib/vpx_image.h"
#include "./lib/vp8cx.h"

// C function pointers
vpx_codec_iface_t *ifaceVP8() {
  return vpx_codec_vp8_cx();
}
vpx_codec_iface_t *ifaceVP9() {
  return vpx_codec_vp9_cx();
}

// C union helpers
void *pktBuf(vpx_codec_cx_pkt_t *pkt) {
  return pkt->data.frame.buf;
}
int pktSz(vpx_codec_cx_pkt_t *pkt) {
  return pkt->data.frame.sz;
}
vpx_codec_frame_flags_t pktFrameFlags(vpx_codec_cx_pkt_t *pkt) {
  return pkt->data.frame.flags;
}

// Alloc helpers
vpx_codec_ctx_t *newCtx() {
  return (vpx_codec_ctx_t*)malloc(sizeof(vpx_codec_ctx_t));
}
vpx_image_t *newImage() {
  return (vpx_image_t*)malloc(sizeof(vpx_image_t));
}

// Wrap encode function to keep Go memory safe
vpx_codec_err_t encode_wrapper(
    vpx_codec_ctx_t* codec, vpx_image_t* raw,
    long t, unsigned long dt, long flags, unsigned long deadline,
    unsigned char *y_ptr, unsigned char *cb_ptr, unsigned char *cr_ptr) {
  raw->planes[0] = y_ptr;
  raw->planes[1] = cb_ptr;
  raw->planes[2] = cr_ptr;
  vpx_codec_err_t ret = vpx_codec_encode(codec, raw, t, dt, flags, deadline);
  raw->planes[0] = raw->planes[1] = raw->planes[2] = 0;
  return ret;
}
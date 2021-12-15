# golang_img_to_neomatrix
golang script to convert any image of the desired side (must be pre-converted) into the neomatrix rgb array

outputs array of the following format for the struct implementation:
{x, y, green, red, blue}

struct pixel {
  int x;
  int y;
  uint8_t green;
  uint8_t red;
  uint8_t blue;
};

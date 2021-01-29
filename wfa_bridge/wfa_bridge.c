#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "gap_affine/affine_wavefront_align.h"

char *load_file(char const *path)
{
  char *buffer = 0;
  long length;
  FILE *f = fopen(path, "rb");

  if (f)
  {
    fseek(f, 0, SEEK_END);
    length = ftell(f);
    fseek(f, 0, SEEK_SET);
    buffer = (char *)malloc((length + 1) * sizeof(char));
    if (buffer)
    {
      fread(buffer, sizeof(char), length, f);
    }
    fclose(f);
  }
  buffer[length] = '\0';

  return buffer;
}

int align(char *reference, char *sequence)
{
  // Allocate MM
  mm_allocator_t *const mm_allocator = mm_allocator_new(BUFFER_SIZE_8M);
  // Set penalties
  affine_penalties_t affine_penalties = {
      .match = 0,
      .mismatch = 4,
      .gap_opening = 6,
      .gap_extension = 2,
  };
  // Init Affine-WFA
  affine_wavefronts_t *affine_wavefronts = affine_wavefronts_new_complete(
      strlen(reference), strlen(sequence), &affine_penalties, NULL, mm_allocator);
  // Align
  affine_wavefronts_align(
      affine_wavefronts, reference, strlen(reference), sequence, strlen(sequence));
  // Display alignment
  const int score = edit_cigar_score_gap_affine(
      &affine_wavefronts->edit_cigar, &affine_penalties);
  fprintf(stderr, "  SCORE COMPUTED %d\t", score);
  edit_cigar_print_pretty(stderr,
                          reference, strlen(reference), sequence, strlen(sequence),
                          &affine_wavefronts->edit_cigar, mm_allocator);
  // Free
  affine_wavefronts_delete(affine_wavefronts);
  mm_allocator_delete(mm_allocator);

  return 0;
}
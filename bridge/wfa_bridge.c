#include <json-c/json.h>
#include "gap_affine/affine_wavefront_align.h"

struct json_object *edit_cigar_json(
    FILE *const stream,
    const char *const pattern,
    const int pattern_length,
    const char *const text,
    const int text_length,
    edit_cigar_t *const edit_cigar,
    mm_allocator_t *const mm_allocator)
{
  // Parameters
  char *const operations = edit_cigar->operations;
  // Allocate alignment buffers
  const int max_buffer_length = text_length + pattern_length + 1;
  char *const pattern_alg = mm_allocator_calloc(mm_allocator, max_buffer_length, char, true);
  char *const ops_alg = mm_allocator_calloc(mm_allocator, max_buffer_length, char, true);
  char *const text_alg = mm_allocator_calloc(mm_allocator, max_buffer_length, char, true);
  // Compute alignment buffers
  int i, alg_pos = 0, pattern_pos = 0, text_pos = 0;
  for (i = edit_cigar->begin_offset; i < edit_cigar->end_offset; ++i)
  {
    switch (operations[i])
    {
    case 'M':
      if (pattern[pattern_pos] != text[text_pos])
      {
        pattern_alg[alg_pos] = pattern[pattern_pos];
        ops_alg[alg_pos] = 'X';
        text_alg[alg_pos++] = text[text_pos];
      }
      else
      {
        pattern_alg[alg_pos] = pattern[pattern_pos];
        ops_alg[alg_pos] = '|';
        text_alg[alg_pos++] = text[text_pos];
      }
      pattern_pos++;
      text_pos++;
      break;
    case 'X':
      if (pattern[pattern_pos] != text[text_pos])
      {
        pattern_alg[alg_pos] = pattern[pattern_pos++];
        ops_alg[alg_pos] = ' ';
        text_alg[alg_pos++] = text[text_pos++];
      }
      else
      {
        pattern_alg[alg_pos] = pattern[pattern_pos++];
        ops_alg[alg_pos] = 'X';
        text_alg[alg_pos++] = text[text_pos++];
      }
      break;
    case 'I':
      pattern_alg[alg_pos] = '-';
      ops_alg[alg_pos] = ' ';
      text_alg[alg_pos++] = text[text_pos++];
      break;
    case 'D':
      pattern_alg[alg_pos] = pattern[pattern_pos++];
      ops_alg[alg_pos] = ' ';
      text_alg[alg_pos++] = '-';
      break;
    default:
      break;
    }
  }
  i = 0;
  while (pattern_pos < pattern_length)
  {
    pattern_alg[alg_pos + i] = pattern[pattern_pos++];
    ops_alg[alg_pos + i] = '?';
    ++i;
  }
  i = 0;
  while (text_pos < text_length)
  {
    text_alg[alg_pos + i] = text[text_pos++];
    ops_alg[alg_pos + i] = '?';
    ++i;
  }
  // edit_cigar_print(stderr, edit_cigar);

  // Free
  mm_allocator_free(mm_allocator, pattern_alg);
  mm_allocator_free(mm_allocator, ops_alg);
  mm_allocator_free(mm_allocator, text_alg);

  struct json_object *jobj;
  jobj = json_object_new_object();
  json_object_object_add(jobj, "pattern_alg", json_object_new_string(pattern_alg));
  json_object_object_add(jobj, "ops_alg", json_object_new_string(ops_alg));
  json_object_object_add(jobj, "text_alg", json_object_new_string(text_alg));

  return jobj;
}

const char *align(char *reference, char *sequence)
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

  struct json_object *jobj;
  jobj = edit_cigar_json(stderr,
                         reference, strlen(reference), sequence, strlen(sequence),
                         &affine_wavefronts->edit_cigar, mm_allocator);
  json_object_object_add(jobj, "score", json_object_new_int(score));

  // Free
  affine_wavefronts_delete(affine_wavefronts);
  mm_allocator_delete(mm_allocator);

  return json_object_to_json_string(jobj);
}
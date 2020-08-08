package cgo

/*
#include "rocksdb/c.h"
#include <string.h>
#include <stdio.h>

void normal_comparator_destroy(void* state) {

}

int key_uint32_comparator_compare(void* state,const char* a, size_t alen,const char* b, size_t blen) {
	  uint32_t a_row_id = a[0] + (a[1] << 8) + (a[2] << 16) + (a[3] << 24);
	  uint32_t b_row_id = b[0] + (b[1] << 8) + (b[2] << 16) + (b[3] << 24);

	  //for (int i=0; i< sizeof(a); i++) {
	  	//	printf("%d ", a[i]);
	  //}

	  //for (int i=0; i< sizeof(b); i++) {
	  	//	printf("%d ", b[i]);
	  //}

      //printf("| %d/%d, %d/%d\n", a_row_id, a[4], b_row_id, b[4]);

	  if (a_row_id == b_row_id) {
	  		int a_column_id = a[4];
	  		int b_column_id = b[4];
	  		if (a_column_id == b_column_id) {
	  			return 0;
	  		} else if (a_column_id < b_column_id) {
	  			return -1;
	  		} else {
	  			return 1;
	  		}
	  } else if (a_row_id < b_row_id) {
	  		return -1;
	  } else {
	  		return 1;
	  }
}

const char* key_uint32_comparator_name(void* state) {
  return "Comparator4uint32";
}

rocksdb_comparator_t* create_key_uint32_comparator() {
	return rocksdb_comparator_create(NULL,
		normal_comparator_destroy,
		key_uint32_comparator_compare,
		key_uint32_comparator_name);
}

int key_uint64_comparator_compare(void* state,const char* a, size_t alen,const char* b, size_t blen) {
	 uint64_t a_row_id = 0;
	 int i = 0;
	 for (i = 0; i < 8; i++)
	 {
	   a_row_id += ((uint64_t) a[i] & 0xFF) << (8 * i);
	 }

	 uint64_t b_row_id = 0;
	 for (i = 0; i < 8; i++)
	 {
	   b_row_id += ((uint64_t) b[i] & 0xFF) << (8 * i);
	 }

     //printf("uint64: %llu/%llu\n", a_row_id, a[8]);

	 if (a_row_id == b_row_id) {
	  		uint64_t a_column_id = a[8];
	  		uint64_t b_column_id = b[8];
	  		if (a_column_id == b_column_id) {
	  			return 0;
	  		} else if (a_column_id < b_column_id) {
	  			return -1;
	  		} else {
	  			return 1;
	  		}
	  } else if (a_row_id < b_row_id) {
	  		return -1;
	  } else {
	  		return 1;
	  }
}

const char* key_uint64_comparator_name(void* state) {
  return "Comparator4uint64";
}

rocksdb_comparator_t* create_key_uint64_comparator() {
	return rocksdb_comparator_create(NULL,
		normal_comparator_destroy,
		key_uint64_comparator_compare,
		key_uint64_comparator_name);
}

int key_int64_comparator_compare(void* state,const char* a, size_t alen,const char* b, size_t blen) {
	 int64_t a_col_val = 0;
	 int i = 0;
	 for (i = 0; i < 8; i++)
	 {
	   a_col_val += ((int64_t) a[i] & 0xFF) << (8 * i);
	 }

	 int64_t b_col_val = 0;
	 for (i = 0; i < 8; i++)
	 {
	   b_col_val += ((int64_t) b[i] & 0xFF) << (8 * i);
	 }

     //printf("uint64: %llu/%llu\n", a_row_id, a[8]);

	 if (a_col_val == b_col_val) {
	 		size_t l_a = alen-8;
	 		size_t l_b = blen-8;

	 		if (l_a < l_b) {
	 			return -1;
	 		} else if (l_a > l_b) {
	 			return 1;
	 		}

	 		for (i=8; i<alen; i++) {
	 			if (((int8_t) a[i]) < ((int8_t) b[i])) {
	 				return -1;
	 			} else if (((int8_t) a[i]) > ((int8_t) b[i])) {
	 				return 1;
	 			}
	 		}

	  		return 0;
	  } else if (a_col_val < b_col_val) {
	  		return -1;
	  } else {
	  		return 1;
	  }
}

const char* key_int64_comparator_name(void* state) {
  return "Comparator4int64";
}

rocksdb_comparator_t* create_key_int64_comparator() {
	return rocksdb_comparator_create(NULL,
		normal_comparator_destroy,
		key_int64_comparator_compare,
		key_int64_comparator_name);
}

*/
import "C"

func NewNativeUInt32Comparator() *C.rocksdb_comparator_t {
	return C.create_key_uint32_comparator()
}

func NewNativeUInt64Comparator() *C.rocksdb_comparator_t {
	return C.create_key_uint64_comparator()
}

func NewNativeInt64Comparator() *C.rocksdb_comparator_t {
	return C.create_key_int64_comparator()
}

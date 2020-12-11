import React from 'react'
import { View, TextInput, StyleSheet } from 'react-native'

const SearchInput = ({ searchText, setSearchText }) => {
  return (
    <View style={styles.searchWrapper}>
      <TextInput
        value={searchText}
        onChangeText={setSearchText}
        placeholder="Input search text"
        style={styles.input}
      />
    </View>
  )
}

export default SearchInput

const styles = StyleSheet.create({
  searchWrapper: {
    width: '100%',
    height: 64,
    paddingVertical: 8,
    paddingHorizontal: 16,
    justifyContent: 'center',
  },
  input: {
    flex: 1,
    borderColor: 'firebrick',
    borderWidth: 2,
    paddingHorizontal: 32,
    borderRadius: 48,
    fontSize: 16,
    backgroundColor: 'azure',
  },
})

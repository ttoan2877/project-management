import React from 'react'
import { StyleSheet, Text, View, TextInput } from 'react-native'

const DataField = ({ label, value, onChange }) => {
  return (
    <View style={styles.wrapper}>
      <Text style={styles.label}>{label}:</Text>
      <TextInput style={styles.input} value={value} onChangeText={onChange} />
    </View>
  )
}

export default DataField

const styles = StyleSheet.create({
  wrapper: {
    marginTop: 16,
  },
  label: {
    fontSize: 12,
    fontWeight: 'bold',
    color: 'azure',
  },
  input: {
    borderColor: 'firebrick',
    borderWidth: 1,
    fontSize: 14,
    padding: 8,
    marginTop: 8,
    borderRadius: 8,
    backgroundColor: 'azure',
  },
})

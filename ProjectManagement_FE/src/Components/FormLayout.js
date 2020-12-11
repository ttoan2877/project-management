import React from 'react'
import { StyleSheet, Text, View } from 'react-native'
import { TouchableOpacity } from 'react-native-gesture-handler'

const FormLayout = ({ title, children, onSubmit, disabled }) => {
  return (
    <View style={styles.container}>
      <Text style={styles.bold}>{title}</Text>
      <View style={styles.formWrapper}>{children}</View>
      <TouchableOpacity
        disabled={disabled}
        onPress={onSubmit}
        style={styles.btn}
      >
        <Text style={styles.label}>Submit</Text>
      </TouchableOpacity>
    </View>
  )
}

export default FormLayout

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginHorizontal: 32,
    marginTop: 16,
    marginBottom: 48,
    borderColor: 'firebrick',
    borderWidth: 4,
    borderRadius: 16,
    padding: 16,
    backgroundColor: 'black',
  },
  formWrapper: {
    flex: 1,
  },
  btn: {
    width: 200,
    height: 40,
    alignSelf: 'center',
    backgroundColor: 'firebrick',
    marginVertical: 16,
    justifyContent: 'center',
    alignItems: 'center',
    borderRadius: 40,
  },
  label: {
    fontSize: 14,
    fontWeight: 'bold',
    color: 'azure',
  },
  bold: {
    fontWeight: 'bold',
    color: 'firebrick',
  },
})
